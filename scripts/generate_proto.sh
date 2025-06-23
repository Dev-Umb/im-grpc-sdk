#!/bin/bash

# Proto文件生成脚本

set -e

# 检查protoc是否安装
if ! command -v protoc &> /dev/null; then
    echo "错误: protoc 未安装"
    echo "请安装 Protocol Buffers compiler:"
    echo "  - Ubuntu/Debian: sudo apt-get install protobuf-compiler"
    echo "  - macOS: brew install protobuf"
    echo "  - 或从 https://github.com/protocolbuffers/protobuf/releases 下载"
    exit 1
fi

# 检查protoc-gen-go是否安装
if ! command -v protoc-gen-go &> /dev/null; then
    echo "错误: protoc-gen-go 未安装"
    echo "请运行: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# 检查protoc-gen-go-grpc是否安装
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "错误: protoc-gen-go-grpc 未安装"
    echo "请运行: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

echo "开始生成 gRPC 代码..."

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# 创建输出目录
OUTPUT_DIR="$PROJECT_DIR/proto/im/v1"
mkdir -p "$OUTPUT_DIR"

# 生成Go代码
protoc \
    --proto_path="$PROJECT_DIR/proto" \
    --go_out="$OUTPUT_DIR" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$OUTPUT_DIR" \
    --go-grpc_opt=paths=source_relative \
    "$PROJECT_DIR/proto/message.proto"

echo "gRPC 代码生成完成!"
echo "生成的文件:"
echo "  - $OUTPUT_DIR/message.pb.go"
echo "  - $OUTPUT_DIR/message_grpc.pb.go"

# 检查生成的文件
if [ -f "$OUTPUT_DIR/message.pb.go" ] && [ -f "$OUTPUT_DIR/message_grpc.pb.go" ]; then
    echo "✅ 所有文件生成成功"
else
    echo "❌ 文件生成失败"
    exit 1
fi 