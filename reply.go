package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// 响应消息结构
type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

// 请求消息结构（与request.go保持一致）
type RequestMessage struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

func main() {
	// 连接到NATS服务器
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("无法连接到NATS服务器: %v", err)
	}
	defer nc.Close()

	// 订阅请求主题
	subject := "service.request"
	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		// 解析请求消息
		var reqMsg RequestMessage
		if err := json.Unmarshal(msg.Data, &reqMsg); err != nil {
			log.Printf("解析请求失败: %v", err)
			return
		}

		// 打印收到的请求
		fmt.Printf("收到请求: %+v\n", reqMsg)

		// 创建响应消息
		respMsg := ResponseMessage{
			Status:  "ok",
			Message: "请求已处理",
			Time:    time.Now().Format(time.RFC3339),
		}

		// 序列化为JSON
		respData, err := json.Marshal(respMsg)
		if err != nil {
			log.Printf("响应序列化失败: %v", err)
			return
		}

		// 发送响应
		if err := msg.Respond(respData); err != nil {
			log.Printf("发送响应失败: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("订阅主题失败: %v", err)
	}

	fmt.Println("等待请求中... (按Ctrl+C退出)")
	// 保持程序运行
	select {}
}
