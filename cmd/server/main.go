package main

import (
	"log"
	"net"

	"github.com/liverlong/rpc-learning/internal/server"
	pb "github.com/liverlong/rpc-learning/pkg/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

func main() {
	// 创建监听器
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建gRPC服务器
	s := grpc.NewServer()

	// 注册用户服务
	userServer := server.NewUserServer()
	pb.RegisterUserServiceServer(s, userServer)

	// 注册反射服务（用于grpcurl等工具）
	reflection.Register(s)

	log.Printf("gRPC server listening on %v", port)
	log.Printf("你可以使用以下命令测试服务:")
	log.Printf("grpcurl -plaintext localhost:50051 list")
	log.Printf("grpcurl -plaintext localhost:50051 user.UserService/ListUsers")

	// 启动服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
