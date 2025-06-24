package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Dev-Umb/im-grpc-sdk/discovery"
	imv1 "github.com/Dev-Umb/im-grpc-sdk/proto/im/v1"
)

// Config 客户端配置
type Config struct {
	// 服务发现配置
	ServiceName  string                     `json:"service_name"`
	Discovery    discovery.ServiceDiscovery `json:"-"`
	LoadBalancer discovery.LoadBalancer     `json:"-"`

	// 连接配置
	ConnectTimeout    time.Duration `json:"connect_timeout"`
	RequestTimeout    time.Duration `json:"request_timeout"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`

	// 重连配置
	MaxRetries    int           `json:"max_retries"`
	RetryInterval time.Duration `json:"retry_interval"`

	// 用户信息
	UserID        string `json:"user_id"`
	DefaultRoomID string `json:"default_room_id"`

	// 回调函数
	OnMessage    func(*imv1.MessageResponse) `json:"-"`
	OnConnect    func()                      `json:"-"`
	OnDisconnect func(error)                 `json:"-"`
	OnError      func(error)                 `json:"-"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		ServiceName:       "im-service",
		ConnectTimeout:    10 * time.Second,
		RequestTimeout:    30 * time.Second,
		HeartbeatInterval: 30 * time.Second,
		MaxRetries:        3,
		RetryInterval:     5 * time.Second,
		LoadBalancer:      discovery.NewRoundRobinBalancer(),
	}
}

// Client IM gRPC客户端
type Client struct {
	config *Config
	conn   *grpc.ClientConn
	client imv1.IMServiceClient
	stream grpc.BidiStreamingClient[imv1.MessageRequest, imv1.MessageResponse]

	// 状态管理
	connected bool
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc

	// 服务发现
	services []*discovery.ServiceInfo

	// 消息处理
	messageCh chan *imv1.MessageRequest

	// 重连
	reconnectCh chan struct{}
}

// NewClient 创建新的IM客户端
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if config.UserID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		config:      config,
		connected:   false,
		ctx:         ctx,
		cancel:      cancel,
		messageCh:   make(chan *imv1.MessageRequest, 100),
		reconnectCh: make(chan struct{}, 1),
	}

	return client, nil
}

// NewClientWithGRPC 使用已有的gRPC客户端创建IM客户端
func NewClientWithGRPC(grpcClient imv1.IMServiceClient, userID string) (*Client, error) {
	if grpcClient == nil {
		return nil, fmt.Errorf("gRPC客户端不能为空")
	}

	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	config := &Config{
		UserID:            userID,
		ConnectTimeout:    10 * time.Second,
		RequestTimeout:    30 * time.Second,
		HeartbeatInterval: 30 * time.Second,
		MaxRetries:        3,
		RetryInterval:     5 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		config:      config,
		client:      grpcClient, // 直接使用传入的gRPC客户端
		connected:   false,
		ctx:         ctx,
		cancel:      cancel,
		messageCh:   make(chan *imv1.MessageRequest, 100),
		reconnectCh: make(chan struct{}, 1),
	}

	return client, nil
}

// NewClientWithGRPCAndConfig 使用已有的gRPC客户端和自定义配置创建IM客户端
func NewClientWithGRPCAndConfig(grpcClient imv1.IMServiceClient, config *Config) (*Client, error) {
	if grpcClient == nil {
		return nil, fmt.Errorf("gRPC客户端不能为空")
	}

	if config == nil {
		return nil, fmt.Errorf("配置不能为空")
	}

	if config.UserID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		config:      config,
		client:      grpcClient, // 直接使用传入的gRPC客户端
		connected:   false,
		ctx:         ctx,
		cancel:      cancel,
		messageCh:   make(chan *imv1.MessageRequest, 100),
		reconnectCh: make(chan struct{}, 1),
	}

	return client, nil
}

// Connect 连接到IM服务
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return fmt.Errorf("客户端已连接")
	}

	// 如果已经有gRPC客户端（通过NewClientWithGRPC创建），跳过连接建立
	if c.client != nil {
		// 直接创建流连接
		if err := c.createStream(); err != nil {
			return fmt.Errorf("创建流连接失败: %v", err)
		}
	} else {
		// 原有的连接建立流程
		// 发现服务
		if err := c.discoverServices(); err != nil {
			return fmt.Errorf("服务发现失败: %v", err)
		}

		// 建立连接
		if err := c.establishConnection(); err != nil {
			return fmt.Errorf("建立连接失败: %v", err)
		}

		// 创建流连接
		if err := c.createStream(); err != nil {
			return fmt.Errorf("创建流连接失败: %v", err)
		}
	}

	c.connected = true

	// 启动后台goroutines
	go c.handleMessages()
	go c.handleHeartbeat()

	// 只有在使用自己管理的连接时才启动重连逻辑
	if c.conn != nil {
		go c.handleReconnect()
	}

	// 监听服务变化（只在有服务发现时）
	if c.config.Discovery != nil {
		go c.watchServices()
	}

	// 触发连接回调
	if c.config.OnConnect != nil {
		c.config.OnConnect()
	}

	return nil
}

// Disconnect 断开连接
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	c.connected = false
	c.cancel()

	if c.stream != nil {
		c.stream.CloseSend()
	}

	// 只关闭自己管理的连接，不关闭注入的gRPC客户端
	if c.conn != nil {
		c.conn.Close()
	}

	return nil
}

// SendTextMessage 发送文本消息
func (c *Client) SendTextMessage(roomID, content string) error {
	return c.SendMessage(&imv1.MessageRequest{
		MessageId: c.generateMessageID(),
		UserId:    c.config.UserID,
		RoomId:    roomID,
		Type:      imv1.MessageType_MESSAGE_TYPE_TEXT,
		Content:   []byte(content),
		Timestamp: timestamppb.New(time.Now()),
	})
}

// SendMessage 发送消息
func (c *Client) SendMessage(msg *imv1.MessageRequest) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return fmt.Errorf("客户端未连接")
	}

	select {
	case c.messageCh <- msg:
		return nil
	case <-time.After(c.config.RequestTimeout):
		return fmt.Errorf("发送消息超时")
	}
}

// JoinRoom 加入房间
func (c *Client) JoinRoom(roomID string, metadata map[string]string) (*imv1.JoinRoomResponse, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return nil, fmt.Errorf("客户端未连接")
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.config.RequestTimeout)
	defer cancel()

	return c.client.JoinRoom(ctx, &imv1.JoinRoomRequest{
		UserId:   c.config.UserID,
		RoomId:   roomID,
		Metadata: metadata,
	})
}

// LeaveRoom 离开房间
func (c *Client) LeaveRoom(roomID string) (*imv1.LeaveRoomResponse, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return nil, fmt.Errorf("客户端未连接")
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.config.RequestTimeout)
	defer cancel()

	return c.client.LeaveRoom(ctx, &imv1.LeaveRoomRequest{
		UserId: c.config.UserID,
		RoomId: roomID,
	})
}

// GetRoomInfo 获取房间信息
func (c *Client) GetRoomInfo(roomID string) (*imv1.GetRoomInfoResponse, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return nil, fmt.Errorf("客户端未连接")
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.config.RequestTimeout)
	defer cancel()

	return c.client.GetRoomInfo(ctx, &imv1.GetRoomInfoRequest{
		RoomId: roomID,
		UserId: c.config.UserID,
	})
}

// UploadAudio 上传音频
func (c *Client) UploadAudio(roomID string, audioData []byte, format string, duration float64) (*imv1.UploadAudioResponse, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return nil, fmt.Errorf("客户端未连接")
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.config.RequestTimeout)
	defer cancel()

	stream, err := c.client.UploadAudio(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建上传流失败: %v", err)
	}

	// 发送元数据
	err = stream.Send(&imv1.UploadAudioRequest{
		Data: &imv1.UploadAudioRequest_Metadata{
			Metadata: &imv1.AudioMetadata{
				UserId:   c.config.UserID,
				RoomId:   roomID,
				Format:   format,
				Size:     int64(len(audioData)),
				Duration: duration,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("发送音频元数据失败: %v", err)
	}

	// 分块发送音频数据
	chunkSize := 1024 * 32 // 32KB chunks
	for i := 0; i < len(audioData); i += chunkSize {
		end := i + chunkSize
		if end > len(audioData) {
			end = len(audioData)
		}

		err = stream.Send(&imv1.UploadAudioRequest{
			Data: &imv1.UploadAudioRequest_Chunk{
				Chunk: audioData[i:end],
			},
		})
		if err != nil {
			return nil, fmt.Errorf("发送音频数据失败: %v", err)
		}
	}

	return stream.CloseAndRecv()
}

// IsConnected 检查连接状态
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// SetServices 设置服务列表（用于直连模式）
func (c *Client) SetServices(services []*discovery.ServiceInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services = services
	if c.config.LoadBalancer != nil {
		c.config.LoadBalancer.Update(services)
	}
}

// discoverServices 发现服务
func (c *Client) discoverServices() error {
	if c.config.Discovery == nil {
		// 如果没有配置服务发现，检查是否已有服务列表
		if len(c.services) == 0 {
			return fmt.Errorf("未配置服务发现且没有可用服务")
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.config.ConnectTimeout)
	defer cancel()

	services, err := c.config.Discovery.Discover(ctx, c.config.ServiceName)
	if err != nil {
		return err
	}

	if len(services) == 0 {
		return fmt.Errorf("未发现可用服务")
	}

	c.services = services
	c.config.LoadBalancer.Update(services)

	return nil
}

// establishConnection 建立gRPC连接
func (c *Client) establishConnection() error {
	if len(c.services) == 0 {
		return fmt.Errorf("没有可用的服务")
	}

	service, err := c.config.LoadBalancer.Select(c.services)
	if err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", service.Address, service.Port)

	ctx, cancel := context.WithTimeout(c.ctx, c.config.ConnectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("连接到 %s 失败: %v", address, err)
	}

	c.conn = conn
	c.client = imv1.NewIMServiceClient(conn)

	return nil
}

// createStream 创建双向流
func (c *Client) createStream() error {
	// 创建带有用户信息的 metadata context
	ctx := c.ctx
	if c.config.UserID != "" && c.config.DefaultRoomID != "" {
		ctx = metadata.AppendToOutgoingContext(c.ctx,
			"user-id", c.config.UserID,
			"room-id", c.config.DefaultRoomID)
	}

	stream, err := c.client.StreamMessages(ctx)
	if err != nil {
		return fmt.Errorf("创建消息流失败: %v", err)
	}

	c.stream = stream

	// 启动接收消息的goroutine
	go c.receiveMessages()

	return nil
}

// handleMessages 处理发送消息
func (c *Client) handleMessages() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.messageCh:
			if err := c.stream.Send(msg); err != nil {
				log.Printf("发送消息失败: %v", err)
				if c.config.OnError != nil {
					c.config.OnError(err)
				}
				// 触发重连
				select {
				case c.reconnectCh <- struct{}{}:
				default:
				}
			}
		}
	}
}

// receiveMessages 接收消息
func (c *Client) receiveMessages() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			msg, err := c.stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Println("服务器关闭了连接")
				} else {
					log.Printf("接收消息失败: %v", err)
				}

				if c.config.OnError != nil {
					c.config.OnError(err)
				}

				// 触发重连
				select {
				case c.reconnectCh <- struct{}{}:
				default:
				}
				return
			}
			if msg.Type == imv1.MessageType_MESSAGE_TYPE_HEARTBEAT {
				continue
			}
			// 处理接收到的消息
			if c.config.OnMessage != nil {
				c.config.OnMessage(msg)
			}
		}
	}
}

// handleHeartbeat 处理心跳
func (c *Client) handleHeartbeat() {
	ticker := time.NewTicker(c.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			heartbeat := &imv1.MessageRequest{
				MessageId: c.generateMessageID(),
				UserId:    c.config.UserID,
				RoomId:    c.config.DefaultRoomID,
				Type:      imv1.MessageType_MESSAGE_TYPE_HEARTBEAT,
				Content:   []byte("ping"),
				Timestamp: timestamppb.New(time.Now()),
			}

			select {
			case c.messageCh <- heartbeat:
			case <-time.After(5 * time.Second):
				log.Println("心跳发送超时")
			}
		}
	}
}

// handleReconnect 处理重连
func (c *Client) handleReconnect() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.reconnectCh:
			c.mu.Lock()
			if c.connected {
				c.connected = false
				if c.config.OnDisconnect != nil {
					c.config.OnDisconnect(fmt.Errorf("连接断开"))
				}

				// 尝试重连
				for i := 0; i < c.config.MaxRetries; i++ {
					log.Printf("尝试重连 (%d/%d)...", i+1, c.config.MaxRetries)

					if err := c.reconnect(); err != nil {
						log.Printf("重连失败: %v", err)
						time.Sleep(c.config.RetryInterval)
						continue
					}

					log.Println("重连成功")
					c.connected = true
					if c.config.OnConnect != nil {
						c.config.OnConnect()
					}
					break
				}
			}
			c.mu.Unlock()
		}
	}
}

// reconnect 重连逻辑
func (c *Client) reconnect() error {
	// 关闭旧流连接
	if c.stream != nil {
		c.stream.CloseSend()
	}

	// 如果使用外部 gRPC 客户端（通过 NewClientWithGRPC 创建），跳过连接重建
	if c.conn == nil {
		// 直接重新创建流
		return c.createStream()
	}

	// 原有的重连逻辑（用于通过服务发现创建的客户端）
	if c.conn != nil {
		c.conn.Close()
	}

	// 重新发现服务
	if err := c.discoverServices(); err != nil {
		return err
	}

	// 重新建立连接
	if err := c.establishConnection(); err != nil {
		return err
	}

	// 重新创建流
	return c.createStream()
}

// watchServices 监听服务变化
func (c *Client) watchServices() {
	serviceCh, err := c.config.Discovery.Watch(c.ctx, c.config.ServiceName)
	if err != nil {
		log.Printf("监听服务变化失败: %v", err)
		return
	}

	for {
		select {
		case <-c.ctx.Done():
			return
		case services := <-serviceCh:
			if services != nil {
				c.services = services
				c.config.LoadBalancer.Update(services)
				log.Printf("服务列表更新: %d个服务", len(services))
			}
		}
	}
}

// generateMessageID 生成消息ID
func (c *Client) generateMessageID() string {
	return fmt.Sprintf("%s_%d", c.config.UserID, time.Now().UnixNano())
}
