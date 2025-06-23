@echo off
setlocal enabledelayedexpansion

REM Proto文件生成脚本 (Windows版本)

echo 开始生成 gRPC 代码...

REM 检查protoc是否安装
where protoc >nul 2>&1
if errorlevel 1 (
    echo 错误: protoc 未安装
    echo 请从 https://github.com/protocolbuffers/protobuf/releases 下载并安装 Protocol Buffers compiler
    exit /b 1
)

REM 检查protoc-gen-go是否安装
where protoc-gen-go >nul 2>&1
if errorlevel 1 (
    echo 错误: protoc-gen-go 未安装
    echo 请运行: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    exit /b 1
)

REM 检查protoc-gen-go-grpc是否安装
where protoc-gen-go-grpc >nul 2>&1
if errorlevel 1 (
    echo 错误: protoc-gen-go-grpc 未安装
    echo 请运行: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    exit /b 1
)

REM 获取脚本所在目录
set SCRIPT_DIR=%~dp0
set PROJECT_DIR=%SCRIPT_DIR%..

REM 创建输出目录
set OUTPUT_DIR=%PROJECT_DIR%\proto\im\v1
if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"

REM 生成Go代码
protoc ^
    --proto_path="%PROJECT_DIR%\proto" ^
    --go_out="%OUTPUT_DIR%" ^
    --go_opt=paths=source_relative ^
    --go-grpc_out="%OUTPUT_DIR%" ^
    --go-grpc_opt=paths=source_relative ^
    "%PROJECT_DIR%\proto\message.proto"

if errorlevel 1 (
    echo ❌ gRPC 代码生成失败
    exit /b 1
)

echo gRPC 代码生成完成!
echo 生成的文件:
echo   - %OUTPUT_DIR%\message.pb.go
echo   - %OUTPUT_DIR%\message_grpc.pb.go

REM 检查生成的文件
if exist "%OUTPUT_DIR%\message.pb.go" if exist "%OUTPUT_DIR%\message_grpc.pb.go" (
    echo ✅ 所有文件生成成功
) else (
    echo ❌ 文件生成失败
    exit /b 1
)

pause 