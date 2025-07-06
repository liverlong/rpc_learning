package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/liverlong/rpc-learning/internal/client"
)

const (
	serverAddr = "localhost:50051"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 3 {
		fmt.Println("用法: go run cmd/chat/main.go <用户ID> <用户名>")
		fmt.Println("示例: go run cmd/chat/main.go 1 张三")
		os.Exit(1)
	}

	// 解析用户ID
	userID, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("无效的用户ID: %v", err)
	}

	username := os.Args[2]

	// 创建客户端
	userClient, err := client.NewUserClient(serverAddr)
	if err != nil {
		log.Fatalf("连接服务器失败: %v", err)
	}
	defer userClient.Close()

	fmt.Printf("=== 聊天客户端 - 用户: %s (ID: %d) ===\n", username, userID)
	fmt.Println("正在连接聊天服务器...")

	// 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 启动聊天
	fmt.Println("开始聊天会话...")
	err = userClient.StartChat(ctx, userID, username)
	if err != nil {
		log.Fatalf("聊天失败: %v", err)
	}

	fmt.Println("聊天会话已结束")
}
