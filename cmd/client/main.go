package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/liverlong/rpc-learning/internal/client"
)

const (
	serverAddr = "localhost:50051"
)

func main() {
	// 创建客户端
	userClient, err := client.NewUserClient(serverAddr)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer userClient.Close()

	fmt.Println("=== gRPC用户服务客户端演示 ===")

	// 1. 创建用户
	fmt.Println("\n1. 创建用户")
	user1, err := userClient.CreateUser("张三", "zhangsan@example.com", 25, "13800138001")
	if err != nil {
		log.Printf("创建用户失败: %v", err)
	} else {
		fmt.Printf("创建的用户: ID=%d, 姓名=%s, 邮箱=%s, 年龄=%d, 电话=%s\n",
			user1.Id, user1.Name, user1.Email, user1.Age, user1.Phone)
	}

	user2, err := userClient.CreateUser("李四", "lisi@example.com", 30, "13800138002")
	if err != nil {
		log.Printf("创建用户失败: %v", err)
	} else {
		fmt.Printf("创建的用户: ID=%d, 姓名=%s, 邮箱=%s, 年龄=%d, 电话=%s\n",
			user2.Id, user2.Name, user2.Email, user2.Age, user2.Phone)
	}

	// 等待一下，确保数据已保存
	time.Sleep(100 * time.Millisecond)

	// 2. 获取用户
	fmt.Println("\n2. 获取用户")
	if user1 != nil {
		getUser, err := userClient.GetUser(user1.Id)
		if err != nil {
			log.Printf("获取用户失败: %v", err)
		} else {
			fmt.Printf("获取的用户: ID=%d, 姓名=%s, 邮箱=%s, 年龄=%d, 电话=%s\n",
				getUser.Id, getUser.Name, getUser.Email, getUser.Age, getUser.Phone)
		}
	}

	// 3. 更新用户
	fmt.Println("\n3. 更新用户")
	if user1 != nil {
		updatedUser, err := userClient.UpdateUser(user1.Id, "张三丰", "zhangsan_new@example.com", 26, "13800138003")
		if err != nil {
			log.Printf("更新用户失败: %v", err)
		} else {
			fmt.Printf("更新的用户: ID=%d, 姓名=%s, 邮箱=%s, 年龄=%d, 电话=%s\n",
				updatedUser.Id, updatedUser.Name, updatedUser.Email, updatedUser.Age, updatedUser.Phone)
		}
	}

	// 4. 列出用户
	fmt.Println("\n4. 列出用户")
	users, total, err := userClient.ListUsers(1, 10)
	if err != nil {
		log.Printf("列出用户失败: %v", err)
	} else {
		fmt.Printf("用户列表 (总数: %d):\n", total)
		for _, user := range users {
			fmt.Printf("  ID=%d, 姓名=%s, 邮箱=%s, 年龄=%d, 电话=%s\n",
				user.Id, user.Name, user.Email, user.Age, user.Phone)
		}
	}

	// 5. 删除用户
	fmt.Println("\n5. 删除用户")
	if user2 != nil {
		err := userClient.DeleteUser(user2.Id)
		if err != nil {
			log.Printf("删除用户失败: %v", err)
		} else {
			fmt.Printf("用户 ID=%d 删除成功\n", user2.Id)
		}
	}

	// 6. 再次列出用户，验证删除
	fmt.Println("\n6. 验证删除后的用户列表")
	users, total, err = userClient.ListUsers(1, 10)
	if err != nil {
		log.Printf("列出用户失败: %v", err)
	} else {
		fmt.Printf("用户列表 (总数: %d):\n", total)
		for _, user := range users {
			fmt.Printf("  ID=%d, 姓名=%s, 邮箱=%s, 年龄=%d, 电话=%s\n",
				user.Id, user.Name, user.Email, user.Age, user.Phone)
		}
	}

	// 7. 测试双向流聊天功能
	fmt.Println("\n7. 测试双向流聊天功能")
	if user1 != nil {
		fmt.Printf("用户 %s 正在加入聊天室...\n", user1.Name)

		// 创建一个带超时的上下文
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// 启动聊天
		err := userClient.StartChat(ctx, user1.Id, user1.Name)
		if err != nil {
			log.Printf("聊天失败: %v", err)
		}
	}

	fmt.Println("\n=== 演示完成 ===")
}
