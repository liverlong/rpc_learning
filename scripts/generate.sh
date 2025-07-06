#!/bin/bash

# 生成protobuf代码的脚本

set -e

echo "开始生成protobuf代码..."

# 检查protoc是否安装
if ! command -v protoc &> /dev/null; then
    echo "错误: protoc 未安装"
    echo "请先安装protoc:"
    echo "  macOS: brew install protobuf"
    echo "  Linux: apt-get install protobuf-compiler"
    exit 1
fi

# 检查protoc-gen-go是否安装
if ! command -v protoc-gen-go &> /dev/null; then
    echo "安装 protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# 检查protoc-gen-go-grpc是否安装
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "安装 protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# 创建输出目录
mkdir -p pkg/pb/user

# 生成Go代码
echo "生成Go代码..."
protoc \
    --go_out=pkg/pb/user \
    --go_opt=paths=source_relative \
    --go-grpc_out=pkg/pb/user \
    --go-grpc_opt=paths=source_relative \
    --proto_path=api/proto \
    api/proto/user.proto

echo "protobuf代码生成完成!"
echo "生成的文件位于: pkg/pb/user/" 