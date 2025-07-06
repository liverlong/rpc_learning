# Makefile for gRPC User Service

.PHONY: help proto build run-server run-client clean test deps install-tools

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  proto         - 生成protobuf代码"
	@echo "  deps          - 下载Go依赖"
	@echo "  build         - 构建服务器和客户端"
	@echo "  run-server    - 运行gRPC服务器"
	@echo "  run-client    - 运行gRPC客户端"
	@echo "  test          - 运行测试"
	@echo "  clean         - 清理构建文件"
	@echo "  install-tools - 安装开发工具"

# 安装开发工具
install-tools:
	@echo "安装protoc工具..."
	@if ! command -v protoc &> /dev/null; then \
		echo "请手动安装protoc:"; \
		echo "  macOS: brew install protobuf"; \
		echo "  Linux: apt-get install protobuf-compiler"; \
		exit 1; \
	fi
	@echo "安装Go protobuf插件..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "安装grpcurl工具..."
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# 生成protobuf代码
proto:
	@echo "生成protobuf代码..."
	@./scripts/generate.sh

# 下载依赖
deps:
	@echo "下载Go依赖..."
	go mod download
	go mod tidy

# 构建项目
build: proto deps
	@echo "构建服务器..."
	go build -o bin/server cmd/server/main.go
	@echo "构建客户端..."
	go build -o bin/client cmd/client/main.go

# 运行服务器
run-server: build
	@echo "启动gRPC服务器..."
	./bin/server

# 运行客户端
run-client: build
	@echo "启动gRPC客户端..."
	./bin/client

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -rf pkg/pb/

# 创建必要的目录
dirs:
	mkdir -p bin
	mkdir -p pkg/pb/user

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "代码检查..."
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint未安装，跳过代码检查"; \
	fi

# 完整构建流程
all: install-tools proto deps build

# 开发环境设置
dev-setup: install-tools proto deps
	@echo "开发环境设置完成!"
	@echo "现在可以运行以下命令:"
	@echo "  make run-server  # 启动服务器"
	@echo "  make run-client  # 启动客户端" 