# IM gRPC SDK Makefile

.PHONY: all clean proto deps test examples install

# 默认目标
all: proto

# 安装依赖
deps:
	@echo "📦 安装依赖..."
	go mod tidy
	go mod download

# 生成proto文件
proto:
	@echo "🔧 生成gRPC代码..."
	@if command -v protoc >/dev/null 2>&1; then \
		chmod +x scripts/generate_proto.sh && \
		./scripts/generate_proto.sh; \
	else \
		echo "❌ protoc 未安装，请先安装 Protocol Buffers compiler"; \
		exit 1; \
	fi

# 安装proto工具
install-proto-tools:
	@echo "🛠️ 安装proto工具..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 测试
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 构建示例
examples: proto
	@echo "🚀 构建示例..."
	go build -o bin/simple_client examples/simple_client.go
	go build -o bin/direct_client examples/direct_client.go

# 清理
clean:
	@echo "🧹 清理文件..."
	rm -rf proto/im/v1/*.pb.go
	rm -rf bin/
	go clean

# 格式化代码
fmt:
	@echo "📝 格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "🔍 代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️ golangci-lint 未安装，使用 go vet"; \
		go vet ./...; \
	fi

# 安装SDK
install:
	@echo "📥 安装SDK..."
	go install

# 创建发布包
release: clean proto
	@echo "📦 创建发布包..."
	mkdir -p release
	tar -czf release/im-grpc-sdk.tar.gz \
		--exclude='.git' \
		--exclude='release' \
		--exclude='bin' \
		.

# 帮助
help:
	@echo "IM gRPC SDK Makefile"
	@echo ""
	@echo "可用目标:"
	@echo "  all                 - 生成proto文件（默认）"
	@echo "  deps                - 安装Go依赖"
	@echo "  proto               - 生成gRPC代码"
	@echo "  install-proto-tools - 安装proto工具"
	@echo "  test                - 运行测试"
	@echo "  examples            - 构建示例程序"
	@echo "  clean               - 清理生成的文件"
	@echo "  fmt                 - 格式化代码"
	@echo "  lint                - 代码检查"
	@echo "  install             - 安装SDK"
	@echo "  release             - 创建发布包"
	@echo "  help                - 显示此帮助" 