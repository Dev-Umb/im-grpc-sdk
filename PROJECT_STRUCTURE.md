# IM gRPC SDK 项目结构

本文档详细说明了 IM gRPC SDK 的项目结构和各个组件的作用。

## 项目概览

```
im_grpc_sdk/
├── 📁 client/                    # 客户端核心实现
│   └── client.go                 # 主要客户端代码
├── 📁 discovery/                 # 服务发现模块
│   ├── interface.go              # 服务发现接口定义
│   ├── consul.go                 # Consul服务发现实现
│   ├── etcd.go                   # ETCD服务发现实现（有依赖问题）
│   └── loadbalancer.go           # 负载均衡器实现
├── 📁 proto/                     # Protocol Buffers定义
│   ├── message.proto             # gRPC服务和消息定义
│   └── 📁 im/v1/                 # 生成的Go代码目录
│       ├── message.pb.go         # 消息类型（需要生成）
│       └── message_grpc.pb.go    # gRPC服务代码（需要生成）
├── 📁 examples/                  # 使用示例
│   ├── simple_client.go          # 基本使用示例（有依赖问题）
│   └── direct_client.go          # 直连模式示例
├── 📁 scripts/                   # 构建和工具脚本
│   ├── generate_proto.sh         # Linux/Mac proto生成脚本
│   └── generate_proto.bat        # Windows proto生成脚本
├── 📄 go.mod                     # Go模块定义
├── 📄 go.sum                     # Go模块依赖锁定
├── 📄 Makefile                   # 构建配置
├── 📄 README.md                  # 主要说明文档
├── 📄 USAGE.md                   # 详细使用指南
└── 📄 PROJECT_STRUCTURE.md       # 本文档
```

## 核心模块详解

### 1. Client 模块 (`client/`)

**文件**: `client.go`

**功能**: 
- IM客户端的核心实现
- 连接管理和自动重连
- 消息发送和接收
- 房间操作（加入、离开、获取信息）
- 音频文件上传
- 心跳保持

**主要结构**:
```go
type Client struct {
    config     *Config                           // 客户端配置
    conn       *grpc.ClientConn                  // gRPC连接
    client     imv1.IMServiceClient              // gRPC客户端
    stream     imv1.IMService_StreamMessagesClient // 双向流
    // ... 其他字段
}
```

**关键方法**:
- `NewClient()` - 创建客户端
- `Connect()` - 连接服务器
- `SendTextMessage()` - 发送文本消息
- `JoinRoom()` - 加入房间
- `UploadAudio()` - 上传音频

### 2. Discovery 模块 (`discovery/`)

#### 2.1 接口定义 (`interface.go`)

定义了服务发现的核心接口：
- `ServiceDiscovery` - 服务发现接口
- `LoadBalancer` - 负载均衡接口
- `HealthChecker` - 健康检查接口

#### 2.2 Consul实现 (`consul.go`)

**功能**:
- Consul服务注册和发现
- 服务健康检查
- 服务变化监听

**使用示例**:
```go
discovery, err := discovery.NewConsulDiscovery("localhost:8500")
```

#### 2.3 ETCD实现 (`etcd.go`)

**状态**: ⚠️ 有依赖问题，需要修复
**功能**:
- ETCD服务注册和发现
- 租约管理
- 服务变化监听

#### 2.4 负载均衡器 (`loadbalancer.go`)

**实现的策略**:
- 轮询 (Round Robin)
- 随机 (Random)
- 加权轮询 (Weighted Round Robin)
- 一致性哈希 (Consistent Hash)

### 3. Proto 模块 (`proto/`)

#### 3.1 消息定义 (`message.proto`)

**定义的服务**:
- `StreamMessages` - 双向流消息
- `SendMessage` - 发送消息
- `JoinRoom` - 加入房间
- `LeaveRoom` - 离开房间
- `GetRoomInfo` - 获取房间信息
- `UploadAudio` - 上传音频
- `GetAudioTranscript` - 获取音频转写
- `HealthCheck` - 健康检查

**消息类型**:
- 文本消息
- 音频消息
- 富文本消息
- 系统消息
- ACK消息
- 房间操作消息
- 心跳消息

#### 3.2 生成的代码 (`im/v1/`)

**需要生成的文件**:
- `message.pb.go` - 消息结构体
- `message_grpc.pb.go` - gRPC服务接口

**生成方法**:
```bash
make proto
# 或
./scripts/generate_proto.sh
```

### 4. Examples 模块 (`examples/`)

#### 4.1 基本示例 (`simple_client.go`)

**状态**: ⚠️ 有依赖问题
**功能**: 展示基本的SDK使用方法

#### 4.2 直连示例 (`direct_client.go`)

**功能**: 
- 不使用服务发现的直连模式
- 完整的消息收发流程
- 房间操作演示

**特点**: 
- 无外部依赖
- 适合快速测试
- 生产环境可用

### 5. Scripts 模块 (`scripts/`)

#### 5.1 Proto生成脚本

**Linux/Mac版本** (`generate_proto.sh`):
- 检查工具安装
- 生成gRPC代码
- 验证生成结果

**Windows版本** (`generate_proto.bat`):
- Windows批处理脚本
- 功能与Linux版本相同

## 构建和使用流程

### 1. 初始化项目

```bash
# 1. 复制SDK到项目目录
cp -r im_grpc_sdk /path/to/your/project/

# 2. 进入SDK目录
cd im_grpc_sdk

# 3. 安装依赖
make deps
```

### 2. 生成Proto代码

```bash
# 安装proto工具（如果未安装）
make install-proto-tools

# 生成代码
make proto
```

### 3. 构建示例

```bash
# 构建所有示例
make examples

# 运行直连示例
./bin/direct_client
```

### 4. 集成到项目

```go
import (
    "github.com/game-im/im-grpc-sdk/client"
    "github.com/game-im/im-grpc-sdk/discovery"
    imv1 "github.com/game-im/im-grpc-sdk/proto/im/v1"
)
```

## 依赖关系

### 核心依赖

```
google.golang.org/grpc v1.59.0
google.golang.org/protobuf v1.31.0
```

### 服务发现依赖

```
github.com/hashicorp/consul/api v1.25.1  # Consul支持
go.etcd.io/etcd/client/v3 v3.5.10       # ETCD支持（有问题）
```

### 工具依赖

```
protoc                    # Protocol Buffers编译器
protoc-gen-go            # Go代码生成器
protoc-gen-go-grpc       # gRPC代码生成器
```

## 已知问题和解决方案

### 1. ETCD依赖问题

**问题**: `go.etcd.io/etcd/client/v3` 导入失败

**解决方案**:
```bash
# 临时解决：注释掉ETCD相关代码
# 或者修复go.mod中的ETCD依赖版本
```

### 2. 示例代码依赖问题

**问题**: 示例代码无法编译

**解决方案**:
```bash
# 先生成proto代码
make proto

# 然后构建示例
make examples
```

### 3. Proto生成失败

**问题**: protoc工具未安装

**解决方案**:
```bash
# Ubuntu/Debian
sudo apt-get install protobuf-compiler

# macOS
brew install protobuf

# 安装Go插件
make install-proto-tools
```

## 使用建议

### 1. 开发环境

- 使用直连模式进行开发和测试
- 启用详细日志进行调试
- 使用 `make examples` 验证功能

### 2. 生产环境

- 使用Consul进行服务发现
- 配置适当的超时和重试参数
- 启用TLS加密连接

### 3. 扩展开发

- 在 `discovery/` 添加新的服务发现实现
- 在 `client/` 扩展客户端功能
- 修改 `proto/message.proto` 添加新的消息类型

## 贡献指南

1. **代码规范**: 遵循Go代码规范
2. **测试**: 添加单元测试
3. **文档**: 更新相关文档
4. **兼容性**: 保持向后兼容

## 版本信息

- **当前版本**: v1.0.0
- **Go版本要求**: 1.21+
- **gRPC版本**: 1.59.0+

## 联系方式

如有问题或建议，请：
1. 提交GitHub Issue
2. 发起Pull Request
3. 联系维护者

---

*最后更新: 2024年* 