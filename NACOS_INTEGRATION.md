# Nacos 集成指南

本文档专门介绍如何将 IM gRPC SDK 与 Nacos 服务发现集成使用。

## 概述

IM gRPC SDK 现在支持直接注入已有的 gRPC 客户端，这使得与 Nacos 等服务发现系统的集成变得非常简单。您无需配置额外的服务发现和负载均衡，直接使用 Nacos 提供的 gRPC 客户端即可。

## 集成优势

✅ **无缝集成**: 直接使用现有的 Nacos 基础设施  
✅ **简化配置**: 无需额外的服务发现配置  
✅ **统一管理**: gRPC 连接由 Nacos SDK 统一管理  
✅ **高可用性**: 利用 Nacos 的服务发现和故障转移  
✅ **性能优化**: 复用已有的连接池和负载均衡  

## 使用方式

### 方式1: 简单集成（推荐）

```go
package main

import (
    "log"
    
    "github.com/Dev-Umb/im-grpc-sdk/client"
    imv1 "github.com/Dev-Umb/im-grpc-sdk/proto/im/v1"
    // "github.com/Dev-Umb/go-pkg/nacos_sdk"
)

func newImServiceClient(conn interface{}) imv1.IMServiceClient {
    return imv1.NewIMServiceClient(conn.(*grpc.ClientConn))
}

func main() {
    // 1. 使用 Nacos 获取 gRPC 客户端
    grpcClient, err := nacos_sdk.GetGRPCClient(
        "im-service",           // 服务名
        "DEFAULT_GROUP",        // Nacos 组
        newImServiceClient,     // 客户端创建函数
    )
    if err != nil {
        log.Fatalf("获取gRPC客户端失败: %v", err)
    }
    
    // 2. 创建 IM 客户端（只需要用户ID）
    imClient, err := client.NewClientWithGRPC(grpcClient, "your_user_id")
    if err != nil {
        log.Fatalf("创建IM客户端失败: %v", err)
    }
    
    // 3. 连接并使用
    if err := imClient.Connect(); err != nil {
        log.Fatalf("连接失败: %v", err)
    }
    defer imClient.Disconnect()
    
    // 4. 使用 IM 功能
    imClient.JoinRoom("room123", nil)
    imClient.SendTextMessage("room123", "Hello from Nacos!")
}
```

### 方式2: 自定义配置

```go
// 1. 获取 Nacos gRPC 客户端（同上）
grpcClient, err := nacos_sdk.GetGRPCClient(
    "im-service", "DEFAULT_GROUP", newImServiceClient)

// 2. 创建自定义配置
config := &client.Config{
    UserID:            "your_user_id",
    DefaultRoomID:     "default_room",
    RequestTimeout:    60 * time.Second,
    HeartbeatInterval: 45 * time.Second,
    
    // 回调函数
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

// 3. 使用自定义配置创建客户端
imClient, err := client.NewClientWithGRPCAndConfig(grpcClient, config)
```

## API 对比

| 功能 | 标准模式 | Nacos集成模式 |
|------|----------|---------------|
| 连接管理 | SDK管理 | Nacos管理 |
| 服务发现 | Consul/ETCD/直连 | Nacos |
| 负载均衡 | SDK内置 | Nacos提供 |
| 健康检查 | SDK实现 | Nacos提供 |
| 重连机制 | SDK处理 | Nacos处理 |
| 配置复杂度 | 较高 | 较低 |

## 配置说明

### 有效配置项

在 Nacos 集成模式下，以下配置项仍然有效：

```go
config := &client.Config{
    UserID:            "user123",        // ✅ 用户ID
    DefaultRoomID:     "room123",        // ✅ 默认房间ID
    RequestTimeout:    60 * time.Second, // ✅ 请求超时
    HeartbeatInterval: 45 * time.Second, // ✅ 心跳间隔
    
    // 回调函数
    OnMessage:         messageHandler,   // ✅ 消息处理
    OnConnect:         connectHandler,   // ✅ 连接成功
    OnDisconnect:      disconnectHandler,// ✅ 连接断开
    OnError:           errorHandler,     // ✅ 错误处理
}
```

### 无效配置项

以下配置项在 Nacos 集成模式下不生效：

```go
config := &client.Config{
    // ❌ 以下配置在 Nacos 模式下无效
    ConnectTimeout:    time.Second,      // 连接由Nacos管理
    MaxRetries:        3,                // 重连由Nacos管理
    RetryInterval:     time.Second,      // 重试由Nacos管理
    Discovery:         discovery,        // 使用Nacos服务发现
    LoadBalancer:      balancer,         // 使用Nacos负载均衡
    ServiceName:       "service",        // 服务名由Nacos管理
}
```

## 完整示例

查看 `examples/nacos_integration.go` 文件获取完整的使用示例。

## 注意事项

1. **连接管理**: 在 Nacos 集成模式下，gRPC 连接完全由 Nacos SDK 管理，IM SDK 不会创建或关闭连接。

2. **错误处理**: 连接相关的错误（如网络断开、服务不可用）会由 Nacos SDK 处理，IM SDK 主要处理业务逻辑错误。

3. **生命周期**: IM 客户端的生命周期独立于 gRPC 连接，您可以多次创建和销毁 IM 客户端而不影响底层连接。

4. **线程安全**: IM 客户端是线程安全的，可以在多个 goroutine 中同时使用。

## 故障排查

### 常见问题

**Q: 创建客户端时提示 "gRPC客户端不能为空"**
```
A: 请确保 nacos_sdk.GetGRPCClient() 成功返回了有效的客户端
```

**Q: 连接成功但无法发送消息**
```
A: 检查用户ID是否正确，以及是否已成功加入房间
```

**Q: 收不到消息回调**
```
A: 确保在创建客户端时正确设置了 OnMessage 回调函数
```

### 调试建议

1. 启用详细日志
2. 检查 Nacos 服务注册状态
3. 验证 gRPC 客户端连接状态
4. 测试基本的房间操作

## 迁移指南

### 从标准模式迁移到 Nacos 模式

1. **移除服务发现配置**:
   ```go
   // 删除这些配置
   config.Discovery = consulDiscovery
   config.LoadBalancer = balancer
   ```

2. **获取 Nacos gRPC 客户端**:
   ```go
   grpcClient, err := nacos_sdk.GetGRPCClient(...)
   ```

3. **更改客户端创建方式**:
   ```go
   // 旧方式
   client, err := client.NewClient(config)
   
   // 新方式
   client, err := client.NewClientWithGRPC(grpcClient, userID)
   ```

## 最佳实践

1. **复用连接**: 在应用程序中复用同一个 gRPC 客户端
2. **错误处理**: 设置适当的错误回调函数
3. **优雅关闭**: 在应用程序退出时调用 `Disconnect()`
4. **监控**: 监控连接状态和消息处理情况

## 版本兼容性

- IM gRPC SDK: v1.1+
- Nacos Go SDK: 兼容主流版本
- Go: 1.21+

---

如有问题，请参考完整文档或提交 Issue。 