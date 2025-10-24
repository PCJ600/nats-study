package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// 请求消息结构
type RequestMessage struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

func main() {
	// 连接到NATS服务器（默认地址）
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("无法连接到NATS服务器: %v", err)
	}
	defer nc.Close()

	// 创建请求消息
	reqMsg := RequestMessage{
		Action:  "ping",
		Message: "Hello from request",
		Time:    time.Now().Format(time.RFC3339),
	}

	// 序列化为JSON
	reqData, err := json.Marshal(reqMsg)
	if err != nil {
		log.Fatalf("JSON序列化失败: %v", err)
	}

	// 发送请求，设置5秒超时
	subject := "service.request"
	resp, err := nc.Request(subject, reqData, 5*time.Second)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	// 打印响应
	fmt.Printf("收到响应: %s\n", resp.Data)
}
