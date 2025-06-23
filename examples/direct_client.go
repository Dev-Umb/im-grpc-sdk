package main

import (
	"log"
	"time"

	"github.com/Dev-Umb/im-grpc-sdk/client"
	"github.com/Dev-Umb/im-grpc-sdk/discovery"
	imv1 "github.com/Dev-Umb/im-grpc-sdk/proto/im/v1"
)

// 直连模式客户端示例 - 不依赖服务发现
func main() {
	log.Println("🚀 启动直连模式IM客户端...")

	// 创建配置
	config := client.DefaultConfig()
	config.UserID = "direct_user_001"
	config.DefaultRoomID = "test_room"

	// 不使用服务发现，直接指定服务地址
	config.Discovery = nil

	// 手动设置服务信息
	services := []*discovery.ServiceInfo{
		{
			ID:      "im-service-direct",
			Name:    "im-service",
			Address: "localhost", // 修改为实际的服务器地址
			Port:    8083,        // 修改为实际的gRPC端口
			Health:  "healthy",
		},
	}

	// 使用轮询负载均衡器
	config.LoadBalancer = discovery.NewRoundRobinBalancer()
	config.LoadBalancer.Update(services)

	// 设置回调函数
	config.OnMessage = func(msg *imv1.MessageResponse) {
		switch msg.Type {
		case imv1.MessageType_MESSAGE_TYPE_TEXT:
			log.Printf("📝 [%s] %s: %s", msg.RoomId, msg.FromUserId, string(msg.Content))
		case imv1.MessageType_MESSAGE_TYPE_SYSTEM:
			log.Printf("🔔 [%s] 系统: %s", msg.RoomId, string(msg.Content))
		case imv1.MessageType_MESSAGE_TYPE_HEARTBEAT:
			// 心跳消息不打印
			return
		default:
			log.Printf("📨 [%s] %s: %v", msg.RoomId, msg.FromUserId, msg.Type)
		}
	}

	config.OnConnect = func() {
		log.Println("✅ 连接成功")
	}

	config.OnDisconnect = func(err error) {
		log.Printf("❌ 连接断开: %v", err)
	}

	config.OnError = func(err error) {
		log.Printf("⚠️ 发生错误: %v", err)
	}

	// 创建客户端
	client, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 设置服务列表（直连模式）
	client.SetServices(services)

	// 连接到服务器
	log.Println("🔗 正在连接服务器...")
	if err := client.Connect(); err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer func() {
		log.Println("👋 断开连接...")
		client.Disconnect()
	}()

	// 等待连接稳定
	time.Sleep(2 * time.Second)

	// 加入房间
	log.Println("🏠 加入房间...")
	joinResp, err := client.JoinRoom("test_room", map[string]string{
		"client":  "direct_example",
		"version": "1.0",
	})
	if err != nil {
		log.Printf("❌ 加入房间失败: %v", err)
	} else {
		log.Printf("✅ 成功加入房间，当前用户数: %d", joinResp.RoomInfo.UserCount)
	}

	// 发送测试消息
	messages := []string{
		"Hello from direct client!",
		"这是一条中文消息",
		"Testing message broadcasting",
		"SDK工作正常",
		"准备结束测试",
	}

	for i, msg := range messages {
		log.Printf("📤 发送消息 %d: %s", i+1, msg)
		if err := client.SendTextMessage("test_room", msg); err != nil {
			log.Printf("❌ 发送消息失败: %v", err)
		} else {
			log.Printf("✅ 消息 %d 发送成功", i+1)
		}

		// 间隔发送
		time.Sleep(3 * time.Second)
	}

	// 获取房间信息
	log.Println("📊 获取房间信息...")
	roomInfo, err := client.GetRoomInfo("test_room")
	if err != nil {
		log.Printf("❌ 获取房间信息失败: %v", err)
	} else {
		log.Printf("📋 房间信息:")
		log.Printf("  - 房间ID: %s", roomInfo.RoomInfo.RoomId)
		log.Printf("  - 用户数: %d", roomInfo.RoomInfo.UserCount)
		log.Printf("  - 消息数: %d", roomInfo.RoomInfo.MessageCount)
		log.Printf("  - 在线用户: %d", len(roomInfo.Users))
	}

	// 保持连接一段时间，接收其他消息
	log.Println("⏳ 保持连接30秒，等待接收消息...")
	time.Sleep(30 * time.Second)

	// 离开房间
	log.Println("🚪 离开房间...")
	_, err = client.LeaveRoom("test_room")
	if err != nil {
		log.Printf("❌ 离开房间失败: %v", err)
	} else {
		log.Println("✅ 成功离开房间")
	}

	log.Println("🎉 直连模式客户端示例完成")
}
