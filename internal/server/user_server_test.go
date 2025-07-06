package server

import (
	"context"
	"testing"

	pb "github.com/liverlong/rpc-learning/pkg/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserServer_CreateUser(t *testing.T) {
	server := NewUserServer()

	tests := []struct {
		name    string
		req     *pb.CreateUserRequest
		wantErr bool
	}{
		{
			name: "valid user",
			req: &pb.CreateUserRequest{
				Name:  "测试用户",
				Email: "test@example.com",
				Age:   25,
				Phone: "13800138000",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: &pb.CreateUserRequest{
				Name:  "",
				Email: "test@example.com",
				Age:   25,
				Phone: "13800138000",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			req: &pb.CreateUserRequest{
				Name:  "测试用户",
				Email: "",
				Age:   25,
				Phone: "13800138000",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.CreateUser(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && resp.User.Name != tt.req.Name {
				t.Errorf("CreateUser() user name = %v, want %v", resp.User.Name, tt.req.Name)
			}
		})
	}
}

func TestUserServer_GetUser(t *testing.T) {
	server := NewUserServer()

	// 先创建一个用户
	createReq := &pb.CreateUserRequest{
		Name:  "测试用户",
		Email: "test@example.com",
		Age:   25,
		Phone: "13800138000",
	}
	createResp, err := server.CreateUser(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	tests := []struct {
		name    string
		req     *pb.GetUserRequest
		wantErr bool
	}{
		{
			name: "existing user",
			req: &pb.GetUserRequest{
				Id: createResp.User.Id,
			},
			wantErr: false,
		},
		{
			name: "non-existing user",
			req: &pb.GetUserRequest{
				Id: 999,
			},
			wantErr: true,
		},
		{
			name: "invalid id",
			req: &pb.GetUserRequest{
				Id: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.GetUser(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && resp.User.Id != tt.req.Id {
				t.Errorf("GetUser() user id = %v, want %v", resp.User.Id, tt.req.Id)
			}
		})
	}
}

func TestUserServer_DeleteUser(t *testing.T) {
	server := NewUserServer()

	// 先创建一个用户
	createReq := &pb.CreateUserRequest{
		Name:  "测试用户",
		Email: "test@example.com",
		Age:   25,
		Phone: "13800138000",
	}
	createResp, err := server.CreateUser(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 删除用户
	deleteReq := &pb.DeleteUserRequest{
		Id: createResp.User.Id,
	}
	_, err = server.DeleteUser(context.Background(), deleteReq)
	if err != nil {
		t.Errorf("DeleteUser() error = %v", err)
	}

	// 验证用户已被删除
	getReq := &pb.GetUserRequest{
		Id: createResp.User.Id,
	}
	_, err = server.GetUser(context.Background(), getReq)
	if err == nil {
		t.Error("Expected error when getting deleted user")
	}
	if status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound error, got %v", status.Code(err))
	}
}

func TestUserServer_ListUsers(t *testing.T) {
	server := NewUserServer()

	// 创建几个用户
	for i := 0; i < 3; i++ {
		createReq := &pb.CreateUserRequest{
			Name:  "测试用户" + string(rune(i+'1')),
			Email: "test" + string(rune(i+'1')) + "@example.com",
			Age:   25 + int32(i),
			Phone: "1380013800" + string(rune(i+'1')),
		}
		_, err := server.CreateUser(context.Background(), createReq)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	// 测试列出用户
	listReq := &pb.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	}
	resp, err := server.ListUsers(context.Background(), listReq)
	if err != nil {
		t.Errorf("ListUsers() error = %v", err)
		return
	}

	if resp.Total != 3 {
		t.Errorf("Expected 3 users, got %d", resp.Total)
	}

	if len(resp.Users) != 3 {
		t.Errorf("Expected 3 users in response, got %d", len(resp.Users))
	}
}
