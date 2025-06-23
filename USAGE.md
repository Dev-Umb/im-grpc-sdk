# IM gRPC SDK 使用指南

本文档提供了 IM gRPC SDK 的详细使用指南和最佳实践。

## 快速开始

### 1. 安装SDK

```bash
# 方式1: 直接复制SDK目录到你的项目中
cp -r im_grpc_sdk /path/to/your/project/

# 方式2: 如果已发布到Git仓库
go get github.com/game-im/im-grpc-sdk
```

### 2. 生成Proto文件

```bash
cd im_grpc_sdk
make proto
# 或者手动执行
./scripts/generate_proto.sh
```

### 3. 使用方式选择

IM gRPC SDK 支持两种主要的使用方式：

#### 方式1: 标准模式（SDK自管理连接）

适用于需要 SDK 自己管理 gRPC 连接和服务发现的场景：

```go
package main

import (
    "log"
    "time"
    
         "github.com/Dev-Umb/im-grpc-sdk/client"
     "github.com/Dev-Umb/im-grpc-sdk/discovery"
     imv1 "github.com/Dev-Umb/im-grpc-sdk/proto/im/v1"
)

func main() {
    // 配置客户端（直连模式）
    config := client.DefaultConfig()
    config.UserID = "user001"
    config.Discovery = nil // 不使用服务发现
    
    // 设置服务地址
    services := []*discovery.ServiceInfo{{
        Address: "localhost",
        Port:    8083,
        Health:  "healthy",
    }}
    config.LoadBalancer = discovery.NewRoundRobinBalancer()
    config.LoadBalancer.Update(services)
    
    // 设置消息回调
    config.OnMessage = func(msg *imv1.MessageResponse) {
        log.Printf("收到消息: %s", string(msg.Content))
    }
    
    // 创建并连接
    client, _ := client.NewClient(config)
    client.Connect()
    defer client.Disconnect()
    
    // 发送消息
    client.SendTextMessage("room001", "Hello World!")
    
    // 保持连接
    time.Sleep(10 * time.Second)
}
```

#### 方式2: Nacos集成模式（使用已有gRPC客户端）

**推荐用于已有Nacos服务发现的项目**，可以直接注入通过Nacos获取的gRPC客户端：

```go
package main

import (
    "log"
    "time"
    
    "github.com/Dev-Umb/im-grpc-sdk/client"
    imv1 "github.com/Dev-Umb/im-grpc-sdk/proto/im/v1"
    // "github.com/Dev-Umb/go-pkg/nacos_sdk" // 您的Nacos SDK
)

func newImServiceClient(conn interface{}) imv1.IMServiceClient {
    // 您的gRPC客户端创建逻辑
    return imv1.NewIMServiceClient(conn.(*grpc.ClientConn))
}

func main() {
    // 使用Nacos获取gRPC客户端
    grpcClient, err := nacos_sdk.GetGRPCClient(
        "im-service",           // 服务名
        "DEFAULT_GROUP",        // Nacos组
        newImServiceClient,     // 客户端创建函数
    )
    if err != nil {
        log.Fatalf("获取gRPC客户端失败: %v", err)
    }
    
    // 简单创建IM客户端（推荐）
    imClient, err := client.NewClientWithGRPC(grpcClient, "user001")
    if err != nil {
        log.Fatal(err)
    }
    
    // 连接（直接使用注入的gRPC客户端，无需额外连接管理）
    if err := imClient.Connect(); err != nil {
        log.Fatal(err)
    }
    defer imClient.Disconnect()
    
    // 发送消息
    imClient.SendTextMessage("room001", "Hello from Nacos!")
    
    // 保持连接
    time.Sleep(10 * time.Second)
}
```

#### Nacos集成的优势

1. **无缝集成**: 直接使用您现有的Nacos服务发现基础设施
2. **统一管理**: gRPC连接由Nacos SDK统一管理，包括负载均衡、健康检查等
3. **简化配置**: 无需额外配置服务发现和负载均衡
4. **高可用**: 利用Nacos的服务发现和故障转移机制
5. **性能优化**: 复用已有的连接池和配置

## 详细配置

### 客户端创建方法对比

| 方法 | 适用场景 | 连接管理 | 服务发现 |
|------|----------|----------|----------|
| `NewClient(config)` | 标准模式 | SDK管理 | 支持多种 |
| `NewClientWithGRPC(grpcClient, userID)` | Nacos集成 | 外部管理 | 由Nacos处理 |
| `NewClientWithGRPCAndConfig(grpcClient, config)` | Nacos集成+自定义 | 外部管理 | 由Nacos处理 |

### 标准模式配置选项

```go
config := &client.Config{
    // === 基本配置 ===
    UserID:        "user123",           // 必填：用户ID
    DefaultRoomID: "default_room",      // 可选：默认房间ID
    ServiceName:   "im-service",        // 服务名称
    
    // === 连接配置 ===
    ConnectTimeout:    10 * time.Second, // 连接超时
    RequestTimeout:    30 * time.Second, // 请求超时
    HeartbeatInterval: 30 * time.Second, // 心跳间隔
    
    // === 重连配置 ===
    MaxRetries:    3,                 // 最大重试次数
    RetryInterval: 5 * time.Second,   // 重试间隔
    
    // === 服务发现和负载均衡 ===
    Discovery:     consulDiscovery,   // 服务发现实例
    LoadBalancer:  loadBalancer,      // 负载均衡器
    
    // === 回调函数 ===
    OnMessage:    messageHandler,     // 消息处理
    OnConnect:    connectHandler,     // 连接成功
    OnDisconnect: disconnectHandler,  // 连接断开
    OnError:      errorHandler,       // 错误处理
}
```

### Nacos集成模式配置

#### 简单配置（推荐）

```go
// 获取Nacos gRPC客户端
grpcClient, err := nacos_sdk.GetGRPCClient(
    "im-service",      // 服务名
    "DEFAULT_GROUP",   // Nacos组
    newImServiceClient,
)

// 直接创建IM客户端
imClient, err := client.NewClientWithGRPC(grpcClient, "user123")
```

#### 高级配置

```go
// 自定义配置
config := &client.Config{
    UserID:            "user123",
    DefaultRoomID:     "default_room",
    RequestTimeout:    60 * time.Second,    // 请求超时
    HeartbeatInterval: 45 * time.Second,    // 心跳间隔
    
    // 注意：以下配置在Nacos模式下不生效
    // ConnectTimeout: 不适用（连接由Nacos管理）
    // MaxRetries:     不适用（重连由Nacos管理）
    // Discovery:      不适用（使用Nacos服务发现）
    // LoadBalancer:   不适用（使用Nacos负载均衡）
    
    // 回调函数仍然有效
    OnMessage: func(msg *imv1.MessageResponse) {
        log.Printf("收到消息: %s", string(msg.Content))
    },
    OnConnect: func() {
        log.Println("IM连接成功")
    },
    OnDisconnect: func(err error) {
        log.Printf("IM连接断开: %v", err)
    },
    OnError: func(err error) {
        log.Printf("IM错误: %v", err)
    },
}

// 使用自定义配置创建客户端
imClient, err := client.NewClientWithGRPCAndConfig(grpcClient, config)
```

## 服务发现配置

### 1. Consul 服务发现

```go
// 创建Consul服务发现
consulDiscovery, err := discovery.NewConsulDiscovery("localhost:8500")
if err != nil {
    log.Fatalf("创建Consul服务发现失败: %v", err)
}

config.Discovery = consulDiscovery
```

### 2. ETCD 服务发现

```go
// 创建ETCD服务发现
etcdDiscovery, err := discovery.NewEtcdDiscovery(
    []string{"localhost:2379"}, 
    "/im-services",
)
if err != nil {
    log.Fatalf("创建ETCD服务发现失败: %v", err)
}

config.Discovery = etcdDiscovery
```

### 3. 直连模式（无服务发现）

```go
// 不使用服务发现
config.Discovery = nil

// 手动设置服务列表
services := []*discovery.ServiceInfo{
    {
        ID:      "im-service-1",
        Address: "10.0.1.100",
        Port:    8083,
        Health:  "healthy",
    },
    {
        ID:      "im-service-2", 
        Address: "10.0.1.101",
        Port:    8083,
        Health:  "healthy",
    },
}

config.LoadBalancer.Update(services)
```

## 负载均衡策略

### 1. 轮询负载均衡

```go
config.LoadBalancer = discovery.NewRoundRobinBalancer()
```

### 2. 随机负载均衡

```go
config.LoadBalancer = discovery.NewRandomBalancer()
```

### 3. 加权轮询负载均衡

```go
config.LoadBalancer = discovery.NewWeightedRoundRobinBalancer()

// 在服务信息中设置权重
services := []*discovery.ServiceInfo{
    {
        Address: "server1",
        Port:    8083,
        Metadata: map[string]string{
            "weight": "3", // 权重为3
        },
    },
    {
        Address: "server2", 
        Port:    8083,
        Metadata: map[string]string{
            "weight": "1", // 权重为1
        },
    },
}
```

### 4. 一致性哈希负载均衡

```go
config.LoadBalancer = discovery.NewConsistentHashBalancer()

// 可以根据用户ID进行哈希
balancer := config.LoadBalancer.(*discovery.ConsistentHashBalancer)
service, err := balancer.SelectByKey(services, userID)
```

## 消息处理

### 消息类型处理

```go
config.OnMessage = func(msg *imv1.MessageResponse) {
    switch msg.Type {
    case imv1.MessageType_MESSAGE_TYPE_TEXT:
        handleTextMessage(msg)
        
    case imv1.MessageType_MESSAGE_TYPE_AUDIO:
        handleAudioMessage(msg)
        
    case imv1.MessageType_MESSAGE_TYPE_RICH_TEXT:
        handleRichTextMessage(msg)
        
    case imv1.MessageType_MESSAGE_TYPE_SYSTEM:
        handleSystemMessage(msg)
        
    case imv1.MessageType_MESSAGE_TYPE_JOIN_ROOM:
        handleUserJoin(msg)
        
    case imv1.MessageType_MESSAGE_TYPE_LEAVE_ROOM:
        handleUserLeave(msg)
        
    default:
        log.Printf("未知消息类型: %v", msg.Type)
    }
    
    // 处理需要确认的消息
    if msg.AckRequired {
        sendAck(msg.MessageId)
    }
}
```

### 发送不同类型的消息

```go
// 1. 文本消息
client.SendTextMessage("room123", "Hello World!")

// 2. 自定义消息
customMsg := &imv1.MessageRequest{
    MessageId: generateMessageID(),
    UserId:    "user123",
    RoomId:    "room123", 
    Type:      imv1.MessageType_MESSAGE_TYPE_RICH_TEXT,
    Content:   []byte(`{"type":"markdown","content":"**Bold Text**"}`),
    Metadata: map[string]string{
        "format": "markdown",
    },
}
client.SendMessage(customMsg)

// 3. 音频消息（先上传音频）
audioData, _ := ioutil.ReadFile("audio.opus")
audioResp, err := client.UploadAudio("room123", audioData, "opus", 30.5)
if err == nil {
    audioMsg := &imv1.MessageRequest{
        MessageId: generateMessageID(),
        UserId:    "user123",
        RoomId:    "room123",
        Type:      imv1.MessageType_MESSAGE_TYPE_AUDIO,
        Content:   []byte(fmt.Sprintf(`{"audio_id":"%s","duration":30.5}`, audioResp.AudioId)),
    }
    client.SendMessage(audioMsg)
}
```

## 房间管理

### 基本房间操作

```go
// 加入房间
joinResp, err := client.JoinRoom("room123", map[string]string{
    "nickname": "张三",
    "role":     "user",
})
if err != nil {
    log.Printf("加入房间失败: %v", err)
} else {
    log.Printf("加入房间成功，当前用户数: %d", joinResp.RoomInfo.UserCount)
}

// 获取房间信息
roomInfo, err := client.GetRoomInfo("room123")
if err != nil {
    log.Printf("获取房间信息失败: %v", err)
} else {
    log.Printf("房间用户数: %d", roomInfo.RoomInfo.UserCount)
    for _, user := range roomInfo.Users {
        log.Printf("用户: %s, 角色: %v", user.UserId, user.Role)
    }
}

// 离开房间
_, err = client.LeaveRoom("room123")
if err != nil {
    log.Printf("离开房间失败: %v", err)
}
```

## 错误处理和重连

### 错误处理最佳实践

```go
config.OnError = func(err error) {
    log.Printf("发生错误: %v", err)
    
    // 可以根据错误类型进行不同处理
    if grpcErr, ok := status.FromError(err); ok {
        switch grpcErr.Code() {
        case codes.Unavailable:
            log.Println("服务不可用，等待重连...")
        case codes.Unauthenticated:
            log.Println("认证失败，请检查用户凭证")
        case codes.PermissionDenied:
            log.Println("权限不足")
        default:
            log.Printf("gRPC错误: %s", grpcErr.Message())
        }
    }
}

config.OnDisconnect = func(err error) {
    log.Printf("连接断开: %v", err)
    // SDK会自动重连，这里可以做一些状态更新
    updateConnectionStatus(false)
}

config.OnConnect = func() {
    log.Println("连接成功")
    updateConnectionStatus(true)
    
    // 重连后可能需要重新加入房间
    rejoinRooms()
}
```

### 自定义重连策略

```go
// 配置重连参数
config.MaxRetries = 10              // 最多重试10次
config.RetryInterval = 2 * time.Second  // 每次重试间隔2秒

// 也可以实现指数退避
func exponentialBackoff(attempt int) time.Duration {
    return time.Duration(math.Pow(2, float64(attempt))) * time.Second
}
```

## 性能优化

### 1. 连接池配置

```go
// 对于高并发场景，可以配置gRPC连接参数
config.ConnectTimeout = 5 * time.Second   // 减少连接超时
config.RequestTimeout = 10 * time.Second  // 减少请求超时
config.HeartbeatInterval = 60 * time.Second // 增加心跳间隔
```

### 2. 消息批处理

```go
// 批量发送消息时，控制发送频率
messages := []string{"msg1", "msg2", "msg3"}
for i, msg := range messages {
    client.SendTextMessage("room123", msg)
    
    // 避免发送过快
    if i < len(messages)-1 {
        time.Sleep(100 * time.Millisecond)
    }
}
```

### 3. 内存优化

```go
// 对于大文件上传，使用流式处理
func uploadLargeAudio(client *client.Client, filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // 分块读取和上传
    buffer := make([]byte, 32*1024) // 32KB chunks
    var audioData []byte
    
    for {
        n, err := file.Read(buffer)
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        audioData = append(audioData, buffer[:n]...)
    }
    
    return client.UploadAudio("room123", audioData, "opus", 0)
}
```

## 安全考虑

### 1. 连接安全

```go
// 生产环境建议使用TLS
conn, err := grpc.Dial(address, 
    grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
)
```

### 2. 消息验证

```go
config.OnMessage = func(msg *imv1.MessageResponse) {
    // 验证消息来源
    if !isValidUser(msg.FromUserId) {
        log.Printf("无效用户消息: %s", msg.FromUserId)
        return
    }
    
    // 验证消息内容
    if len(msg.Content) > maxMessageSize {
        log.Printf("消息过长: %d bytes", len(msg.Content))
        return
    }
    
    // 处理消息
    handleMessage(msg)
}
```

## 故障排查

### 常见问题

1. **连接失败**
   ```
   检查服务器地址和端口是否正确
   检查网络连接
   检查防火墙设置
   ```

2. **服务发现失败**
   ```
   检查Consul/ETCD是否运行
   检查服务注册是否成功
   检查网络连接
   ```

3. **消息发送失败**
   ```
   检查是否已连接
   检查房间是否存在
   检查用户权限
   ```

### 调试模式

```go
// 启用详细日志
config.OnError = func(err error) {
    log.Printf("详细错误信息: %+v", err)
}

// 监控连接状态
go func() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if client.IsConnected() {
                log.Println("连接状态: 已连接")
            } else {
                log.Println("连接状态: 未连接")
            }
        }
    }
}()
```

## 部署建议

### 1. 生产环境配置

```go
productionConfig := &client.Config{
    ConnectTimeout:    15 * time.Second,
    RequestTimeout:    60 * time.Second,
    HeartbeatInterval: 120 * time.Second,
    MaxRetries:        5,
    RetryInterval:     10 * time.Second,
    // 使用生产环境的服务发现
    Discovery: consulDiscovery,
}
```

### 2. 监控和指标

```go
// 添加指标收集
var (
    messagesSent     int64
    messagesReceived int64
    connectionErrors int64
)

config.OnMessage = func(msg *imv1.MessageResponse) {
    atomic.AddInt64(&messagesReceived, 1)
    handleMessage(msg)
}

config.OnError = func(err error) {
    atomic.AddInt64(&connectionErrors, 1)
    handleError(err)
}
```

## 总结

IM gRPC SDK 提供了完整的即时通讯功能，支持多种部署模式和配置选项。通过合理配置和使用最佳实践，可以构建稳定、高性能的IM应用。

更多示例和详细信息，请参考 `examples/` 目录下的示例代码。 