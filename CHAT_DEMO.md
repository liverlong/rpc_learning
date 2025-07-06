# gRPC 双向流聊天功能演示

本项目已经添加了双向流聊天功能，演示了如何使用 gRPC 的双向流进行实时通信。

## 功能特性

- **双向流通信**: 客户端和服务端可以同时发送和接收消息
- **实时聊天**: 多个用户可以同时在线聊天
- **用户管理**: 自动管理用户加入和离开
- **消息广播**: 消息会实时广播给所有在线用户
- **在线用户统计**: 实时显示在线用户数量

## 架构设计

### Protocol Buffer 定义

```protobuf
// 聊天消息
message ChatMessage {
  int64 user_id = 1;
  string username = 2;
  string content = 3;
  int64 timestamp = 4;
  string message_type = 5; // text, join, leave, system
}

// 聊天请求
message ChatRequest {
  int64 user_id = 1;
  string username = 2;
  string content = 3;
  string action = 4; // join, leave, message
}

// 聊天响应
message ChatResponse {
  ChatMessage message = 1;
  string status = 2;
  int32 online_users = 3;
}

// 双向流聊天接口
rpc Chat(stream ChatRequest) returns (stream ChatResponse);
```

### 服务端实现

- **连接管理**: 维护所有活跃的聊天连接
- **消息广播**: 将消息广播给所有在线用户
- **用户状态**: 跟踪用户的加入和离开
- **并发安全**: 使用互斥锁保护共享数据

### 客户端实现

- **双向通信**: 同时处理发送和接收消息
- **异步处理**: 使用 goroutine 处理消息接收
- **自动重连**: 处理连接断开和重连逻辑

## 使用方法

### 1. 启动服务器

```bash
# 构建并启动服务器
make run-server
```

服务器将在 `localhost:50051` 上启动。

### 2. 启动聊天客户端

打开多个终端窗口，分别启动不同的聊天客户端：

```bash
# 终端1 - 用户张三
make run-chat USER_ID=1 USERNAME=张三

# 终端2 - 用户李四
make run-chat USER_ID=2 USERNAME=李四

# 终端3 - 用户王五
make run-chat USER_ID=3 USERNAME=王五
```

### 3. 观察聊天过程

每个客户端会：
1. 自动加入聊天室
2. 发送几条测试消息
3. 接收其他用户的消息
4. 在超时后自动离开

## 消息类型

- **system**: 系统消息（欢迎消息等）
- **join**: 用户加入消息
- **leave**: 用户离开消息
- **text**: 普通聊天消息

## 测试场景

### 场景1: 单用户聊天
```bash
make run-chat USER_ID=1 USERNAME=张三
```

### 场景2: 多用户聊天
```bash
# 在不同终端中同时运行
make run-chat USER_ID=1 USERNAME=张三 &
make run-chat USER_ID=2 USERNAME=李四 &
make run-chat USER_ID=3 USERNAME=王五 &
```

### 场景3: 用户动态加入和离开
```bash
# 先启动一个用户
make run-chat USER_ID=1 USERNAME=张三

# 等待几秒后启动另一个用户
make run-chat USER_ID=2 USERNAME=李四

# 观察用户加入和离开的消息广播
```

## 日志输出示例

### 服务端日志
```
2024/01/01 10:00:00 Chat stream started
2024/01/01 10:00:00 Received chat request: {UserId:1 Username:张三 Content: Action:join}
2024/01/01 10:00:02 Received chat request: {UserId:1 Username:张三 Content:大家好！ Action:message}
2024/01/01 10:00:03 Chat stream ended
```

### 客户端日志
```
2024/01/01 10:00:00 聊天连接已建立，用户: 张三 (ID: 1)
2024/01/01 10:00:00 [系统] 欢迎 张三 加入聊天室！ (10:00:00) - 在线用户: 1
2024/01/01 10:00:01 [加入] 李四 加入了聊天室 (10:00:01) - 在线用户: 2
2024/01/01 10:00:02 [李四] 大家好！ (10:00:02)
2024/01/01 10:00:03 [张三] 我是 张三 (10:00:03)
```

## 扩展功能

可以基于此实现进一步扩展：

1. **私聊功能**: 支持用户间的私人消息
2. **聊天室管理**: 支持多个聊天室
3. **消息持久化**: 将消息保存到数据库
4. **用户认证**: 添加用户身份验证
5. **消息过滤**: 添加敏感词过滤
6. **文件传输**: 支持文件和图片传输
7. **消息历史**: 支持查看历史消息

## 技术要点

1. **双向流**: 使用 `stream` 关键字定义双向流接口
2. **并发控制**: 使用互斥锁保护共享数据结构
3. **goroutine**: 使用 goroutine 处理并发的发送和接收
4. **上下文管理**: 使用 context 控制连接生命周期
5. **错误处理**: 妥善处理网络错误和连接断开

这个双向流聊天功能展示了 gRPC 在实时通信场景中的强大能力，为构建更复杂的实时应用奠定了基础。 