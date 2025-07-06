# gRPC用户服务Demo

这是一个基于Go和Protocol Buffers的gRPC服务示例项目，实现了完整的用户管理CRUD操作和双向流聊天功能。

## 项目结构

```
rpc_learning/
├── api/proto/              # Protocol Buffer定义文件
│   └── user.proto         # 用户服务定义
├── cmd/                   # 应用程序入口点
│   ├── server/           # gRPC服务器
│   │   └── main.go
│   ├── client/           # gRPC客户端
│   │   └── main.go
│   └── chat/             # 聊天客户端
│       └── main.go
├── internal/             # 内部包
│   ├── server/          # 服务器实现
│   │   └── user_server.go
│   └── client/          # 客户端实现
│       └── user_client.go
├── pkg/pb/              # 生成的Protocol Buffer代码
│   └── user/           # 用户服务相关代码
├── scripts/             # 构建脚本
│   ├── generate.sh     # 生成protobuf代码
│   └── demo_chat.sh    # 聊天功能演示
├── go.mod              # Go模块文件
├── go.sum              # Go依赖校验文件
├── Makefile            # 构建配置
├── README.md           # 项目说明
└── CHAT_DEMO.md        # 聊天功能详细说明
```

## 功能特性

### 基础CRUD操作
- ✅ 创建用户 (CreateUser)
- ✅ 获取用户 (GetUser)
- ✅ 更新用户 (UpdateUser)
- ✅ 删除用户 (DeleteUser)
- ✅ 列出用户 (ListUsers) - 支持分页

### 双向流聊天功能 🆕
- ✅ 双向流通信 (Chat)
- ✅ 实时消息传递
- ✅ 多用户同时在线
- ✅ 用户加入/离开通知
- ✅ 在线用户统计
- ✅ 消息广播

### 技术特性
- ✅ 完整的错误处理
- ✅ 参数验证
- ✅ 并发安全
- ✅ 超时控制

## 前置要求

- Go 1.21+
- Protocol Buffers编译器 (protoc)
- Make工具

### 安装Protocol Buffers编译器

**macOS:**
```bash
brew install protobuf
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get install protobuf-compiler
```

**Windows:**
下载预编译的二进制文件：https://github.com/protocolbuffers/protobuf/releases

## 快速开始

### 1. 克隆项目并进入目录
```bash
git clone <repository-url>
cd rpc_learning
```

### 2. 安装开发工具和依赖
```bash
make install-tools
make deps
```

### 3. 生成Protocol Buffer代码
```bash
make proto
```

### 4. 构建项目
```bash
make build
```

### 5. 运行服务器
```bash
make run-server
```

服务器将在 `localhost:50051` 上启动。

### 6. 运行客户端（新终端）
```bash
make run-client
```

客户端将连接到服务器并演示所有功能。

### 7. 体验双向流聊天功能 🚀
```bash
# 自动演示（推荐）
make demo-chat

# 或者手动测试
# 终端1: make run-server
# 终端2: make run-chat USER_ID=1 USERNAME=张三
# 终端3: make run-chat USER_ID=2 USERNAME=李四
```

## 使用方法

### 使用Makefile命令

```bash
make help          # 显示所有可用命令
make install-tools # 安装开发工具
make proto         # 生成protobuf代码
make deps          # 下载依赖
make build         # 构建项目
make run-server    # 运行服务器
make run-client    # 运行客户端
make run-chat      # 运行聊天客户端 (需要参数)
make demo-chat     # 演示双向流聊天功能
make test          # 运行测试
make clean         # 清理构建文件
make fmt           # 格式化代码
make lint          # 代码检查
```

### 手动构建和运行

**生成protobuf代码:**
```bash
./scripts/generate.sh
```

**构建服务器:**
```bash
go build -o bin/server cmd/server/main.go
```

**构建客户端:**
```bash
go build -o bin/client cmd/client/main.go
```

**构建聊天客户端:**
```bash
go build -o bin/chat cmd/chat/main.go
```

**运行服务器:**
```bash
./bin/server
```

**运行客户端:**
```bash
./bin/client
```

**运行聊天客户端:**
```bash
./bin/chat 1 张三
```

## API接口

### 用户服务 (UserService)

#### 基础CRUD操作

#### 1. CreateUser - 创建用户
```protobuf
rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
```

#### 2. GetUser - 获取用户
```protobuf
rpc GetUser(GetUserRequest) returns (GetUserResponse);
```

#### 3. UpdateUser - 更新用户
```protobuf
rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
```

#### 4. DeleteUser - 删除用户
```protobuf
rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
```

#### 5. ListUsers - 列出用户
```protobuf
rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
```

#### 双向流接口

#### 6. Chat - 双向流聊天 🆕
```protobuf
rpc Chat(stream ChatRequest) returns (stream ChatResponse);
```

支持的操作：
- `join`: 加入聊天室
- `message`: 发送消息
- `leave`: 离开聊天室

## 双向流聊天功能

### 快速体验

1. **自动演示**（推荐新手）：
   ```bash
   make demo-chat
   ```

2. **手动测试**（推荐深入了解）：
   ```bash
   # 终端1: 启动服务器
   make run-server
   
   # 终端2: 用户张三加入聊天
   make run-chat USER_ID=1 USERNAME=张三
   
   # 终端3: 用户李四加入聊天
   make run-chat USER_ID=2 USERNAME=李四
   ```

### 功能特性

- **实时通信**: 基于gRPC双向流的实时消息传递
- **多用户支持**: 支持多个用户同时在线聊天
- **消息广播**: 消息自动广播给所有在线用户
- **状态管理**: 自动处理用户加入和离开
- **在线统计**: 实时显示在线用户数量
- **消息类型**: 支持系统消息、用户消息、加入/离开通知

### 详细说明

查看 [CHAT_DEMO.md](CHAT_DEMO.md) 获取双向流聊天功能的详细说明和技术实现。

## 测试工具

### 使用grpcurl测试

安装grpcurl:
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

测试命令:
```bash
# 列出所有服务
grpcurl -plaintext localhost:50051 list

# 列出用户服务的方法
grpcurl -plaintext localhost:50051 list user.UserService

# 创建用户
grpcurl -plaintext -d '{"name":"测试用户","email":"test@example.com","age":25,"phone":"13800138000"}' localhost:50051 user.UserService/CreateUser

# 获取用户
grpcurl -plaintext -d '{"id":1}' localhost:50051 user.UserService/GetUser

# 列出用户
grpcurl -plaintext -d '{"page":1,"page_size":10}' localhost:50051 user.UserService/ListUsers
```

## 项目特点

### 1. 标准的Go项目结构
- 遵循Go项目布局标准
- 清晰的目录分层
- 合理的包组织

### 2. 完整的gRPC实现
- 定义清晰的Protocol Buffer接口
- 完整的服务器端实现
- 易用的客户端封装
- 支持gRPC反射
- **双向流通信支持** 🆕

### 3. 生产级别的代码质量
- 完整的错误处理
- 参数验证
- 并发安全
- 超时控制
- 日志记录

### 4. 开发友好
- 自动化构建脚本
- 详细的文档
- 示例代码
- 测试工具支持
- **交互式演示** 🆕

## 扩展建议

1. **数据持久化**: 集成数据库（如PostgreSQL、MySQL）
2. **认证授权**: 添加JWT认证和权限控制
3. **监控指标**: 集成Prometheus监控
4. **链路追踪**: 添加OpenTelemetry支持
5. **配置管理**: 使用Viper进行配置管理
6. **单元测试**: 添加完整的单元测试覆盖
7. **Docker化**: 添加Docker支持
8. **CI/CD**: 集成GitHub Actions
9. **聊天增强**: 私聊、聊天室、消息持久化 🆕
10. **文件传输**: 基于流的文件上传下载 🆕

## 故障排除

### 常见问题

1. **protoc未找到**
   - 确保已安装Protocol Buffers编译器
   - 检查PATH环境变量

2. **生成代码失败**
   - 确保protoc-gen-go和protoc-gen-go-grpc已安装
   - 检查proto文件语法

3. **连接失败**
   - 确保服务器已启动
   - 检查端口是否被占用

4. **导入错误**
   - 运行 `go mod tidy` 更新依赖
   - 确保生成的pb文件存在

5. **聊天功能问题** 🆕
   - 确保服务器支持双向流
   - 检查客户端连接状态
   - 查看服务器日志

## 贡献指南

1. Fork项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License 