# IM gRPC SDK

一个功能完整的即时通讯 gRPC 客户端 SDK，支持服务发现、负载均衡、自动重连和双向流通信。

## 特性

- 🚀 **高性能**: 基于 gRPC 协议，支持双向流通信
- 🔍 **服务发现**: 支持 Consul、ETCD 等服务发现机制
- ⚖️ **负载均衡**: 内置多种负载均衡策略（轮询、随机、加权轮询、一致性哈希）
- 🔄 **自动重连**: 智能重连机制，保证连接稳定性
- 💬 **多消息类型**: 支持文本、音频、富文本、系统消息等多种类型
- 🏠 **房间管理**: 完整的房间加入、离开、信息获取功能
- 📁 **文件上传**: 支持流式音频文件上传
- 🔧 **易于集成**: 简单的 API 设计，快速集成到现有项目

## 安装

```bash
go get github.com/Dev-Umb/im-grpc-sdk
```

## 快速开始

### 1. 生成 Proto 文件

首先需要生成 gRPC 代码：

```bash
# Linux/macOS
chmod +x scripts/generate_proto.sh
./scripts/generate_proto.sh

# Windows
scripts\generate_proto.bat
```

### 2. 使用方式

SDK 支持两种使用方式：

#### 方式1: 基本使用（SDK自管理连接）

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
    // 创建服务发现（可选）
    consulDiscovery, err := discovery.NewConsulDiscovery("localhost:8500")
    if err != nil {
        log.Fatalf("创建服务发现失败: %v", err)
    }

    // 配置客户端
    config := client.DefaultConfig()
    config.UserID = "user123"
    config.DefaultRoomID = "room456"
    config.Discovery = consulDiscovery
    config.LoadBalancer = discovery.NewRoundRobinBalancer()

    // 设置消息回调
    config.OnMessage = func(msg *imv1.MessageResponse) {
        log.Printf("收到消息: %s", string(msg.Content))
    }

    // 创建并连接客户端
    client, err := client.NewClient(config)
    if err != nil {
        log.Fatalf("创建客户端失败: %v", err)
    }

    if err := client.Connect(); err != nil {
        log.Fatalf("连接失败: %v", err)
    }
    defer client.Disconnect()

    // 加入房间
    _, err = client.JoinRoom("room456", nil)
    if err != nil {
        log.Printf("加入房间失败: %v", err)
    }

    // 发送消息
    err = client.SendTextMessage("room456", "Hello World!")
    if err != nil {
        log.Printf("发送消息失败: %v", err)
    }

    // 保持连接
    time.Sleep(30 * time.Second)
}
```

#### 方式2: 使用已有gRPC客户端（推荐用于Nacos集成）

如果您已经有 gRPC 客户端（比如通过 Nacos 服务发现获取的），可以直接注入使用：

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
    
    // 方式2.1: 简单创建（推荐）
    imClient, err := client.NewClientWithGRPC(grpcClient, "user123")
    if err != nil {
        log.Fatal(err)
    }
    
    // 方式2.2: 使用自定义配置
    /*
    config := &client.Config{
        UserID:            "user123",
        DefaultRoomID:     "room456",
        RequestTimeout:    60 * time.Second,
        HeartbeatInterval: 45 * time.Second,
        OnMessage: func(msg *imv1.MessageResponse) {
            log.Printf("收到消息: %s", string(msg.Content))
        },
        OnConnect: func() {
            log.Println("连接成功")
        },
        OnDisconnect: func(err error) {
            log.Printf("连接断开: %v", err)
        },
    }
    
    imClient, err := client.NewClientWithGRPCAndConfig(grpcClient, config)
    if err != nil {
        log.Fatal(err)
    }
    */
    
    // 连接（不会创建新的gRPC连接，直接使用注入的客户端）
    if err := imClient.Connect(); err != nil {
        log.Fatal(err)
    }
    defer imClient.Disconnect()
    
    // 加入房间
    _, err = imClient.JoinRoom("room456", map[string]string{
        "source": "nacos_integration",
    })
    if err != nil {
        log.Printf("加入房间失败: %v", err)
    }
    
    // 发送消息
    err = imClient.SendTextMessage("room456", "Hello from Nacos!")
    if err != nil {
        log.Printf("发送消息失败: %v", err)
    }
    
    // 保持连接
    time.Sleep(30 * time.Second)
}
```

### 3. 无服务发现模式

如果不使用服务发现，可以直接连接到固定地址：

```go
// 创建直连客户端
config := client.DefaultConfig()
config.UserID = "user123"
config.Discovery = nil // 不使用服务发现

// 手动设置服务信息
services := []*discovery.ServiceInfo{
    {
        ID:      "im-service-1",
        Name:    "im-service",
        Address: "localhost",
        Port:    8083,
        Health:  "healthy",
    },
}

config.LoadBalancer = discovery.NewRoundRobinBalancer()
config.LoadBalancer.Update(services)
```

## API 文档

### 客户端配置

```go
type Config struct {
    // 服务发现配置
    ServiceName     string
    Discovery       discovery.ServiceDiscovery
    LoadBalancer    discovery.LoadBalancer
    
    // 连接配置
    ConnectTimeout    time.Duration  // 连接超时时间
    RequestTimeout    time.Duration  // 请求超时时间
    HeartbeatInterval time.Duration  // 心跳间隔
    
    // 重连配置
    MaxRetries      int           // 最大重试次数
    RetryInterval   time.Duration // 重试间隔
    
    // 用户信息
    UserID          string        // 用户ID（必填）
    DefaultRoomID   string        // 默认房间ID
    
    // 回调函数
    OnMessage       func(*imv1.MessageResponse) // 消息接收回调
    OnConnect       func()                      // 连接成功回调
    OnDisconnect    func(error)                 // 连接断开回调
    OnError         func(error)                 // 错误回调
}
```

### 客户端创建方法

```go
// 标准方式：创建客户端（SDK自管理连接）
func NewClient(config *Config) (*Client, error)

// Nacos集成：使用已有gRPC客户端创建（简单版本）
func NewClientWithGRPC(grpcClient imv1.IMServiceClient, userID string) (*Client, error)

// Nacos集成：使用已有gRPC客户端和自定义配置创建
func NewClientWithGRPCAndConfig(grpcClient imv1.IMServiceClient, config *Config) (*Client, error)
```

### 主要方法

#### 连接管理

```go
// 连接到服务器
func (c *Client) Connect() error

// 断开连接
func (c *Client) Disconnect() error

// 检查连接状态
func (c *Client) IsConnected() bool
```

#### 消息发送

```go
// 发送文本消息
func (c *Client) SendTextMessage(roomID, content string) error

// 发送自定义消息
func (c *Client) SendMessage(msg *imv1.MessageRequest) error
```

#### 房间操作

```go
// 加入房间
func (c *Client) JoinRoom(roomID string, metadata map[string]string) (*imv1.JoinRoomResponse, error)

// 离开房间
func (c *Client) LeaveRoom(roomID string) (*imv1.LeaveRoomResponse, error)

// 获取房间信息
func (c *Client) GetRoomInfo(roomID string) (*imv1.GetRoomInfoResponse, error)
```

#### 文件上传

```go
// 上传音频文件
func (c *Client) UploadAudio(roomID string, audioData []byte, format string, duration float64) (*imv1.UploadAudioResponse, error)
```

## 服务发现

SDK 支持多种服务发现机制：

### Consul

```go
consulDiscovery, err := discovery.NewConsulDiscovery("localhost:8500")
if err != nil {
    log.Fatalf("创建Consul服务发现失败: %v", err)
}

config.Discovery = consulDiscovery
```

### ETCD

```go
etcdDiscovery, err := discovery.NewEtcdDiscovery([]string{"localhost:2379"}, "/services")
if err != nil {
    log.Fatalf("创建ETCD服务发现失败: %v", err)
}

config.Discovery = etcdDiscovery
```

## 负载均衡

SDK 内置多种负载均衡策略：

### 轮询负载均衡

```go
config.LoadBalancer = discovery.NewRoundRobinBalancer()
```

### 随机负载均衡

```go
config.LoadBalancer = discovery.NewRandomBalancer()
```

### 加权轮询负载均衡

```go
config.LoadBalancer = discovery.NewWeightedRoundRobinBalancer()
```

### 一致性哈希负载均衡

```go
config.LoadBalancer = discovery.NewConsistentHashBalancer()
```

## 消息类型

SDK 支持以下消息类型：

- `MESSAGE_TYPE_TEXT`: 文本消息
- `MESSAGE_TYPE_AUDIO`: 音频消息
- `MESSAGE_TYPE_RICH_TEXT`: 富文本消息
- `MESSAGE_TYPE_SYSTEM`: 系统消息
- `MESSAGE_TYPE_ACK`: 确认消息
- `MESSAGE_TYPE_JOIN_ROOM`: 加入房间消息
- `MESSAGE_TYPE_LEAVE_ROOM`: 离开房间消息
- `MESSAGE_TYPE_HEARTBEAT`: 心跳消息

## 错误处理

SDK 提供完善的错误处理机制：

```go
config.OnError = func(err error) {
    log.Printf("发生错误: %v", err)
    // 可以在这里实现自定义错误处理逻辑
}

config.OnDisconnect = func(err error) {
    log.Printf("连接断开: %v", err)
    // 连接断开时的处理逻辑
    // SDK 会自动尝试重连
}
```

## 配置选项

### 默认配置

```go
&Config{
    ServiceName:       "im-service",
    ConnectTimeout:    10 * time.Second,
    RequestTimeout:    30 * time.Second,
    HeartbeatInterval: 30 * time.Second,
    MaxRetries:        3,
    RetryInterval:     5 * time.Second,
    LoadBalancer:      discovery.NewRoundRobinBalancer(),
}
```

### 自定义配置

```go
config := client.DefaultConfig()
config.ConnectTimeout = 15 * time.Second
config.RequestTimeout = 60 * time.Second
config.HeartbeatInterval = 60 * time.Second
config.MaxRetries = 5
config.RetryInterval = 10 * time.Second
```

## 完整示例

查看 `examples/` 目录下的完整示例：

- `simple_client.go`: 基本使用示例
- `advanced_client.go`: 高级功能示例
- `batch_client.go`: 批量操作示例

## 依赖要求

- Go 1.21+
- Protocol Buffers 3.0+
- gRPC-Go 1.59.0+

## 安装依赖工具

### Protocol Buffers Compiler

```bash
# Ubuntu/Debian
sudo apt-get install protobuf-compiler

# macOS
brew install protobuf

# 或从 https://github.com/protocolbuffers/protobuf/releases 下载
```

### Go Proto 插件

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 项目结构

```
im_grpc_sdk/
├── client/           # 客户端实现
├── discovery/        # 服务发现实现
├── proto/           # Proto 文件和生成的代码
├── examples/        # 使用示例
├── scripts/         # 构建脚本
├── go.mod          # Go 模块文件
└── README.md       # 说明文档
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 更新日志

### v1.0.0
- 初始版本发布
- 支持基本的 IM 功能
- 集成服务发现和负载均衡
- 支持自动重连
- 完整的文档和示例

## 支持

如有问题，请提交 Issue 或联系维护者。

## 新功能：gRPC Metadata 支持

### 自动用户信息传递

从 v1.2.0 开始，SDK 支持通过 gRPC metadata 自动传递用户信息，无需在连接后发送初始消息。

#### 优势

1. **简化连接流程**：不需要发送初始连接消息
2. **更快的连接建立**：直接通过 metadata 传递用户信息
3. **向后兼容**：仍支持原有的初始消息方式
4. **更符合 gRPC 最佳实践**：使用 metadata 传递元信息

#### 使用方法

```go
// 创建 gRPC 客户端
grpcClient := imv1.NewIMServiceClient(conn)

// 配置 IM 客户端
config := &client.Config{
    UserID:        "user123",
    DefaultRoomID: "room456",
    // ... 其他配置
}

// 创建 IM 客户端（自动使用 metadata 传递用户信息）
imClient, err := client.NewClientWithGRPCAndConfig(grpcClient, config)
if err != nil {
    log.Fatalf("创建IM客户端失败: %v", err)
}

// 连接到服务器（会自动通过 metadata 传递 userID 和 roomID）
err = imClient.Connect()
if err != nil {
    log.Fatalf("连接失败: %v", err)
}

// 现在可以直接发送消息，无需发送初始连接消息
err = imClient.SendTextMessage("room456", "Hello!")
```

#### 技术实现

- 客户端在创建流连接时自动添加 `user-id` 和 `room-id` metadata
- 服务端优先从 metadata 读取用户信息
- 如果 metadata 中没有用户信息，自动回退到原有的初始消息方式