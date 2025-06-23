package main

import (
	"log"
	"time"

	"github.com/Dev-Umb/im-grpc-sdk/client"
	imv1 "github.com/Dev-Umb/im-grpc-sdk/proto/im/v1"
	// 假设您有这个包，根据实际情况调整
	// "github.com/Dev-Umb/go-pkg/nacos_sdk"
)

// 模拟您的 newImServiceClient 函数
func newImServiceClient(conn interface{}) imv1.IMServiceClient {
	// 这里应该是您实际的客户端创建逻辑
	// 返回真实的 IMServiceClient
	return nil // 占位符，实际使用时请替换
}

func main() {
	log.Println("🚀 启动Nacos集成IM客户端示例...")

	// 使用Nacos获取gRPC客户端（根据您的实际代码调整）
	/*
		config := &YourConfig{
			SubLoginServerName: "im-service",
			NacosGroup:        "DEFAULT_GROUP",
		}

		grpcClient, err := nacos_sdk.GetGRPCClient(
			config.SubLoginServerName,
			config.NacosGroup,
			newImServiceClient,
		)
		if err != nil {
			log.Fatalf("获取gRPC客户端失败: %v", err)
		}
	*/

	// 方式1: 简单初始化（只需要用户ID）
	grpcClient := newImServiceClient(nil) // 替换为您的实际客户端
	if grpcClient == nil {
		log.Println("⚠️ 示例中使用模拟客户端，请替换为实际的gRPC客户端")
		log.Println("💡 实际使用示例:")
		log.Println("   grpcClient, err := nacos_sdk.GetGRPCClient(serviceName, group, newImServiceClient)")
		log.Println("   imClient, err := client.NewClientWithGRPC(grpcClient, userID)")
		return
	}

	imClient, err := client.NewClientWithGRPC(grpcClient, "nacos_user_123")
	if err != nil {
		log.Fatalf("创建IM客户端失败: %v", err)
	}

	// 方式2: 使用自定义配置（注释掉的完整示例）
	/*
		config := &client.Config{
			UserID:            "nacos_user_123",
			DefaultRoomID:     "nacos_room",
			RequestTimeout:    60 * time.Second,
			HeartbeatInterval: 45 * time.Second,
			OnMessage: func(msg *imv1.MessageResponse) {
				log.Printf("📨 收到消息: [%s] %s: %s",
					msg.RoomId, msg.FromUserId, string(msg.Content))
			},
			OnConnect: func() {
				log.Println("✅ IM客户端连接成功")
			},
			OnDisconnect: func(err error) {
				log.Printf("❌ IM客户端连接断开: %v", err)
			},
			OnError: func(err error) {
				log.Printf("⚠️ IM客户端发生错误: %v", err)
			},
		}

		imClient, err := client.NewClientWithGRPCAndConfig(grpcClient, config)
		if err != nil {
			log.Fatalf("创建IM客户端失败: %v", err)
		}
	*/

	// 注意：消息回调应该在创建客户端时通过Config设置
	// 这里只是演示，实际使用时请在Config中设置OnMessage回调

	// 连接到服务器
	log.Println("🔗 正在连接IM服务...")
	if err := imClient.Connect(); err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer func() {
		log.Println("👋 断开IM连接...")
		imClient.Disconnect()
	}()

	// 等待连接稳定
	time.Sleep(2 * time.Second)

	// 加入房间
	log.Println("🏠 加入房间...")
	joinResp, err := imClient.JoinRoom("nacos_room", map[string]string{
		"client":  "nacos_example",
		"version": "1.0",
		"source":  "nacos_discovery",
	})
	if err != nil {
		log.Printf("❌ 加入房间失败: %v", err)
	} else {
		log.Printf("✅ 成功加入房间，当前用户数: %d", joinResp.RoomInfo.UserCount)
	}

	// 发送测试消息
	messages := []string{
		"Hello from Nacos integrated client!",
		"Nacos服务发现集成测试",
		"gRPC客户端注入成功",
		"消息发送正常",
		"准备结束测试",
	}

	for i, msg := range messages {
		log.Printf("📤 发送消息 %d: %s", i+1, msg)
		if err := imClient.SendTextMessage("nacos_room", msg); err != nil {
			log.Printf("❌ 发送消息失败: %v", err)
		} else {
			log.Printf("✅ 消息 %d 发送成功", i+1)
		}

		// 间隔发送
		time.Sleep(3 * time.Second)
	}

	// 获取房间信息
	log.Println("📊 获取房间信息...")
	roomInfo, err := imClient.GetRoomInfo("nacos_room")
	if err != nil {
		log.Printf("❌ 获取房间信息失败: %v", err)
	} else {
		log.Printf("📋 房间信息:")
		log.Printf("  - 房间ID: %s", roomInfo.RoomInfo.RoomId)
		log.Printf("  - 用户数: %d", roomInfo.RoomInfo.UserCount)
		log.Printf("  - 消息数: %d", roomInfo.RoomInfo.MessageCount)
	}

	// 保持连接一段时间，接收其他消息
	log.Println("⏳ 保持连接30秒，等待接收消息...")
	time.Sleep(30 * time.Second)

	// 离开房间
	log.Println("🚪 离开房间...")
	_, err = imClient.LeaveRoom("nacos_room")
	if err != nil {
		log.Printf("❌ 离开房间失败: %v", err)
	} else {
		log.Println("✅ 成功离开房间")
	}

	log.Println("🎉 Nacos集成IM客户端示例完成")
}
