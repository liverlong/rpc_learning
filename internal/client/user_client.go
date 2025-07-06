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
