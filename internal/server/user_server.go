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

// ChatClient 聊天客户端信息
type ChatClient struct {
	UserID   int64
	Username string
	Stream   pb.UserService_ChatServer
}

// UserServer 用户服务服务器
type UserServer struct {
	pb.UnimplementedUserServiceServer
	users       map[int64]*pb.User
	nextID      int64
	mu          sync.RWMutex
	chatClients map[int64]*ChatClient
	chatMu      sync.RWMutex
}

// NewUserServer 创建新的用户服务服务器
func NewUserServer() *UserServer {
	return &UserServer{
		users:       make(map[int64]*pb.User),
		nextID:      1,
		chatClients: make(map[int64]*ChatClient),
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

// Chat 双向流聊天接口
func (s *UserServer) Chat(stream pb.UserService_ChatServer) error {
	log.Printf("Chat stream started")

	var client *ChatClient
	defer func() {
		// 清理客户端连接
		if client != nil {
			s.chatMu.Lock()
			delete(s.chatClients, client.UserID)
			s.chatMu.Unlock()

			// 广播用户离开消息
			s.broadcastMessage(&pb.ChatMessage{
				UserId:      client.UserID,
				Username:    client.Username,
				Content:     fmt.Sprintf("%s 离开了聊天室", client.Username),
				Timestamp:   time.Now().Unix(),
				MessageType: "leave",
			}, client.UserID)
		}
		log.Printf("Chat stream ended")
	}()

	for {
		// 接收客户端消息
		req, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return err
		}

		log.Printf("Received chat request: %+v", req)

		switch req.Action {
		case "join":
			// 用户加入聊天室
			client = &ChatClient{
				UserID:   req.UserId,
				Username: req.Username,
				Stream:   stream,
			}

			s.chatMu.Lock()
			s.chatClients[req.UserId] = client
			onlineUsers := len(s.chatClients)
			s.chatMu.Unlock()

			// 发送加入确认
			joinResponse := &pb.ChatResponse{
				Message: &pb.ChatMessage{
					UserId:      0,
					Username:    "系统",
					Content:     fmt.Sprintf("欢迎 %s 加入聊天室！", req.Username),
					Timestamp:   time.Now().Unix(),
					MessageType: "system",
				},
				Status:      "joined",
				OnlineUsers: int32(onlineUsers),
			}

			if err := stream.Send(joinResponse); err != nil {
				log.Printf("Error sending join response: %v", err)
				return err
			}

			// 广播用户加入消息
			s.broadcastMessage(&pb.ChatMessage{
				UserId:      req.UserId,
				Username:    req.Username,
				Content:     fmt.Sprintf("%s 加入了聊天室", req.Username),
				Timestamp:   time.Now().Unix(),
				MessageType: "join",
			}, req.UserId)

		case "message":
			// 处理聊天消息
			if client == nil {
				// 用户未加入聊天室
				errorResponse := &pb.ChatResponse{
					Message: &pb.ChatMessage{
						UserId:      0,
						Username:    "系统",
						Content:     "请先加入聊天室",
						Timestamp:   time.Now().Unix(),
						MessageType: "system",
					},
					Status:      "error",
					OnlineUsers: 0,
				}

				if err := stream.Send(errorResponse); err != nil {
					log.Printf("Error sending error response: %v", err)
					return err
				}
				continue
			}

			// 广播用户消息
			s.broadcastMessage(&pb.ChatMessage{
				UserId:      req.UserId,
				Username:    req.Username,
				Content:     req.Content,
				Timestamp:   time.Now().Unix(),
				MessageType: "text",
			}, 0) // 0表示广播给所有用户

		case "leave":
			// 用户主动离开
			if client != nil {
				s.chatMu.Lock()
				delete(s.chatClients, client.UserID)
				s.chatMu.Unlock()

				// 广播用户离开消息
				s.broadcastMessage(&pb.ChatMessage{
					UserId:      client.UserID,
					Username:    client.Username,
					Content:     fmt.Sprintf("%s 离开了聊天室", client.Username),
					Timestamp:   time.Now().Unix(),
					MessageType: "leave",
				}, client.UserID)

				client = nil
			}
			return nil

		default:
			log.Printf("Unknown action: %s", req.Action)
		}
	}
}

// broadcastMessage 广播消息给所有在线用户
func (s *UserServer) broadcastMessage(message *pb.ChatMessage, excludeUserID int64) {
	s.chatMu.RLock()
	defer s.chatMu.RUnlock()

	response := &pb.ChatResponse{
		Message:     message,
		Status:      "broadcast",
		OnlineUsers: int32(len(s.chatClients)),
	}

	for userID, client := range s.chatClients {
		if excludeUserID != 0 && userID == excludeUserID {
			continue // 跳过指定用户
		}

		if err := client.Stream.Send(response); err != nil {
			log.Printf("Error broadcasting to user %d: %v", userID, err)
			// 删除断开连接的客户端
			delete(s.chatClients, userID)
		}
	}
}
