package main

import (
	"log"
	"time"

	"github.com/Dev-Umb/im-grpc-sdk/client"
	"github.com/Dev-Umb/im-grpc-sdk/discovery"
	imv1 "github.com/Dev-Umb/im-grpc-sdk/proto/im/v1"
)

func main() {
	// 创建Consul服务发现
	consulDiscovery, err := discovery.NewConsulDiscovery("localhost:8500")
	if err != nil {
		log.Fatalf("创建Consul服务发现失败: %v", err)
	}

	// 配置客户端
	config := client.DefaultConfig()
	config.UserID = "user123"
	config.DefaultRoomID = "room456"
	config.Discovery = consulDiscovery
	config.LoadBalancer = discovery.NewRoundRobinBalancer()

	// 设置回调函数
	config.OnMessage = func(msg *imv1.MessageResponse) {
		log.Printf("收到消息: 类型=%v, 来自=%s, 内容=%s",
			msg.Type, msg.FromUserId, string(msg.Content))
	}

	config.OnConnect = func() {
		log.Println("连接成功")
	}

	config.OnDisconnect = func(err error) {
		log.Printf("连接断开: %v", err)
	}

	config.OnError = func(err error) {
		log.Printf("发生错误: %v", err)
	}

	// 创建客户端
	client, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 连接到服务器
	if err := client.Connect(); err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Disconnect()

	// 加入房间
	joinResp, err := client.JoinRoom("room456", map[string]string{
		"client": "simple_example",
	})
	if err != nil {
		log.Printf("加入房间失败: %v", err)
	} else {
		log.Printf("加入房间成功: %s", joinResp.Status.Message)
	}

	// 发送文本消息
	if err := client.SendTextMessage("room456", "Hello from SDK!"); err != nil {
		log.Printf("发送消息失败: %v", err)
	}

	// 获取房间信息
	roomInfo, err := client.GetRoomInfo("room456")
	if err != nil {
		log.Printf("获取房间信息失败: %v", err)
	} else {
		log.Printf("房间信息: 用户数=%d", roomInfo.RoomInfo.UserCount)
	}

	// 等待一段时间接收消息
	time.Sleep(30 * time.Second)

	log.Println("客户端示例结束")
}
