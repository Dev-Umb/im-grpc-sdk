package main

import (
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Dev-Umb/im-grpc-sdk/client"
	imv1 "github.com/Dev-Umb/im-grpc-sdk/proto/im/v1"
)

func main() {
	// 直接连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	grpcClient := imv1.NewIMServiceClient(conn)

	// 配置 IM 客户端
	config := &client.Config{
		UserID:            "user123",
		DefaultRoomID:     "room456",
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
			log.Printf("❌ IM客户端错误: %v", err)
		},
	}

	// 创建 IM 客户端（使用 metadata 方式）
	imClient, err := client.NewClientWithGRPCAndConfig(grpcClient, config)
	if err != nil {
		log.Fatalf("创建IM客户端失败: %v", err)
	}

	// 连接到服务器
	err = imClient.Connect()
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer imClient.Disconnect()

	log.Println("客户端已连接，使用 metadata 方式传递用户信息")

	// 等待一下连接建立
	time.Sleep(2 * time.Second)

	// 发送测试消息
	err = imClient.SendTextMessage("room456", "Hello from metadata client!")
	if err != nil {
		log.Printf("发送消息失败: %v", err)
	}

	// 发送第二条消息
	time.Sleep(1 * time.Second)
	err = imClient.SendTextMessage("room456", "This message was sent without initial connection message!")
	if err != nil {
		log.Printf("发送消息失败: %v", err)
	}

	// 等待接收消息
	log.Println("等待接收消息...")
	time.Sleep(30 * time.Second)

	log.Println("测试完成")
}
