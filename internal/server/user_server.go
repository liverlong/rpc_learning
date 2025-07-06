package server

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/liverlong/rpc-learning/pkg/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer 用户服务服务器
type UserServer struct {
	pb.UnimplementedUserServiceServer
	users  map[int64]*pb.User
	nextID int64
	mu     sync.RWMutex
}

// NewUserServer 创建新的用户服务服务器
func NewUserServer() *UserServer {
	return &UserServer{
		users:  make(map[int64]*pb.User),
		nextID: 1,
	}
}

// CreateUser 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("CreateUser called with: %+v", req)

	// 验证请求参数
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "用户名不能为空")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "邮箱不能为空")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查邮箱是否已存在
	for _, user := range s.users {
		if user.Email == req.Email {
			return nil, status.Error(codes.AlreadyExists, "邮箱已存在")
		}
	}

	// 创建新用户
	now := time.Now().Unix()
	user := &pb.User{
		Id:        s.nextID,
		Name:      req.Name,
		Email:     req.Email,
		Age:       req.Age,
		Phone:     req.Phone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.users[s.nextID] = user
	s.nextID++

	return &pb.CreateUserResponse{
		User:    user,
		Message: "用户创建成功",
	}, nil
}

// GetUser 获取用户
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("GetUser called with: %+v", req)

	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户ID必须大于0")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	return &pb.GetUserResponse{
		User:    user,
		Message: "获取用户成功",
	}, nil
}

// UpdateUser 更新用户
func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	log.Printf("UpdateUser called with: %+v", req)

	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户ID必须大于0")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != "" && req.Email != user.Email {
		for _, u := range s.users {
			if u.Id != req.Id && u.Email == req.Email {
				return nil, status.Error(codes.AlreadyExists, "邮箱已被其他用户使用")
			}
		}
	}

	// 更新用户信息
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Age > 0 {
		user.Age = req.Age
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	user.UpdatedAt = time.Now().Unix()

	return &pb.UpdateUserResponse{
		User:    user,
		Message: "用户更新成功",
	}, nil
}

// DeleteUser 删除用户
func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	log.Printf("DeleteUser called with: %+v", req)

	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户ID必须大于0")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.users[req.Id]
	if !exists {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	delete(s.users, req.Id)

	return &pb.DeleteUserResponse{
		Message: "用户删除成功",
	}, nil
}

// ListUsers 列出用户
func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	log.Printf("ListUsers called with: %+v", req)

	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// 获取所有用户
	allUsers := make([]*pb.User, 0, len(s.users))
	for _, user := range s.users {
		allUsers = append(allUsers, user)
	}

	total := int32(len(allUsers))

	// 计算分页
	start := (page - 1) * pageSize
	end := start + pageSize

	var users []*pb.User
	if start < total {
		if end > total {
			end = total
		}
		users = allUsers[start:end]
	} else {
		users = []*pb.User{}
	}

	return &pb.ListUsersResponse{
		Users:   users,
		Total:   total,
		Message: fmt.Sprintf("获取用户列表成功，共%d个用户", total),
	}, nil
}
