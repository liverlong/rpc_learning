#!/bin/bash

# 双向流聊天功能演示脚本

set -e

echo "=== gRPC 双向流聊天功能演示 ==="
echo

# 检查是否已构建
if [ ! -f "bin/server" ] || [ ! -f "bin/chat" ]; then
    echo "构建项目..."
    make build
    echo
fi

echo "1. 启动服务器..."
echo "   在后台启动 gRPC 服务器"
./bin/server &
SERVER_PID=$!

# 等待服务器启动
sleep 2

echo
echo "2. 启动多个聊天客户端..."
echo "   将启动3个客户端进行聊天演示"
echo

# 启动聊天客户端（在后台运行）
echo "   启动用户: 张三 (ID: 1)"
./bin/chat 1 张三 &
CLIENT1_PID=$!

sleep 1

echo "   启动用户: 李四 (ID: 2)"
./bin/chat 2 李四 &
CLIENT2_PID=$!

sleep 1

echo "   启动用户: 王五 (ID: 3)"
./bin/chat 3 王五 &
CLIENT3_PID=$!

echo
echo "3. 观察聊天过程..."
echo "   客户端将自动发送消息并接收其他用户的消息"
echo "   请观察终端输出中的聊天消息"
echo

# 等待聊天完成
sleep 15

echo
echo "4. 清理进程..."

# 终止所有进程
kill $CLIENT1_PID $CLIENT2_PID $CLIENT3_PID $SERVER_PID 2>/dev/null || true

# 等待进程结束
sleep 2

echo
echo "=== 演示完成 ==="
echo
echo "如果想手动测试，可以运行："
echo "  终端1: make run-server"
echo "  终端2: make run-chat USER_ID=1 USERNAME=张三"
echo "  终端3: make run-chat USER_ID=2 USERNAME=李四"
echo 