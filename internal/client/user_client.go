package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/liverlong/rpc-learning/pkg/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserClient 用户服务客户端
type UserClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

// NewUserClient 创建新的用户服务客户端
func NewUserClient(serverAddr string) (*UserClient, error) {
	// 建立连接
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	// 创建客户端
	client := pb.NewUserServiceClient(conn)

	return &UserClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close 关闭连接
func (c *UserClient) Close() error {
	return c.conn.Close()
}

// CreateUser 创建用户
func (c *UserClient) CreateUser(name, email string, age int32, phone string) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.CreateUserRequest{
		Name:  name,
		Email: email,
		Age:   age,
		Phone: phone,
	}

	resp, err := c.client.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	log.Printf("创建用户成功: %s", resp.Message)
	return resp.User, nil
}

// GetUser 获取用户
func (c *UserClient) GetUser(id int64) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetUserRequest{Id: id}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	log.Printf("获取用户成功: %s", resp.Message)
	return resp.User, nil
}

// UpdateUser 更新用户
func (c *UserClient) UpdateUser(id int64, name, email string, age int32, phone string) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.UpdateUserRequest{
		Id:    id,
		Name:  name,
		Email: email,
		Age:   age,
		Phone: phone,
	}

	resp, err := c.client.UpdateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	log.Printf("更新用户成功: %s", resp.Message)
	return resp.User, nil
}

// DeleteUser 删除用户
func (c *UserClient) DeleteUser(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.DeleteUserRequest{Id: id}

	resp, err := c.client.DeleteUser(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	log.Printf("删除用户成功: %s", resp.Message)
	return nil
}

// ListUsers 列出用户
func (c *UserClient) ListUsers(page, pageSize int32) ([]*pb.User, int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.ListUsers(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %v", err)
	}

	log.Printf("获取用户列表成功: %s", resp.Message)
	return resp.Users, resp.Total, nil
}

// StartChat 启动聊天功能
func (c *UserClient) StartChat(ctx context.Context, userID int64, username string) error {
	// 创建双向流
	stream, err := c.client.Chat(ctx)
	if err != nil {
		return fmt.Errorf("failed to start chat: %v", err)
	}

	log.Printf("聊天连接已建立，用户: %s (ID: %d)", username, userID)

	// 发送加入聊天室的请求
	joinReq := &pb.ChatRequest{
		UserId:   userID,
		Username: username,
		Action:   "join",
	}

	if err := stream.Send(joinReq); err != nil {
		return fmt.Errorf("failed to join chat: %v", err)
	}

	// 启动接收消息的goroutine
	go c.receiveMessages(stream)

	// 启动发送消息的goroutine
	go c.sendMessages(stream, userID, username)

	// 等待上下文取消
	<-ctx.Done()

	// 发送离开消息
	leaveReq := &pb.ChatRequest{
		UserId:   userID,
		Username: username,
		Action:   "leave",
	}

	if err := stream.Send(leaveReq); err != nil {
		log.Printf("Failed to send leave message: %v", err)
	}

	// 关闭流
	if err := stream.CloseSend(); err != nil {
		log.Printf("Failed to close send stream: %v", err)
	}

	log.Printf("聊天连接已关闭")
	return nil
}

// receiveMessages 接收服务器消息
func (c *UserClient) receiveMessages(stream pb.UserService_ChatClient) {
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}

		if resp.Message != nil {
			timestamp := time.Unix(resp.Message.Timestamp, 0)
			switch resp.Message.MessageType {
			case "system":
				log.Printf("[系统] %s (%s) - 在线用户: %d",
					resp.Message.Content,
					timestamp.Format("15:04:05"),
					resp.OnlineUsers)
			case "join":
				log.Printf("[加入] %s (%s) - 在线用户: %d",
					resp.Message.Content,
					timestamp.Format("15:04:05"),
					resp.OnlineUsers)
			case "leave":
				log.Printf("[离开] %s (%s) - 在线用户: %d",
					resp.Message.Content,
					timestamp.Format("15:04:05"),
					resp.OnlineUsers)
			case "text":
				log.Printf("[%s] %s (%s)",
					resp.Message.Username,
					resp.Message.Content,
					timestamp.Format("15:04:05"))
			}
		}
	}
}

// sendMessages 发送消息（这里简化处理，实际应用中可以从标准输入读取）
func (c *UserClient) sendMessages(stream pb.UserService_ChatClient, userID int64, username string) {
	// 模拟发送一些测试消息
	messages := []string{
		"大家好！",
		"我是 " + username,
		"很高兴认识大家！",
	}

	for i, content := range messages {
		time.Sleep(time.Duration(i+2) * time.Second) // 间隔发送

		req := &pb.ChatRequest{
			UserId:   userID,
			Username: username,
			Content:  content,
			Action:   "message",
		}

		if err := stream.Send(req); err != nil {
			log.Printf("Failed to send message: %v", err)
			return
		}
	}
}

// SendChatMessage 发送单条聊天消息（用于交互式聊天）
func (c *UserClient) SendChatMessage(stream pb.UserService_ChatClient, userID int64, username, content string) error {
	req := &pb.ChatRequest{
		UserId:   userID,
		Username: username,
		Content:  content,
		Action:   "message",
	}

	return stream.Send(req)
}
