package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// 与pub.go保持一致的消息结构
type PubMessage struct {
	Topic   string `json:"topic"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

func main() {
	// 连接到NATS服务器
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("无法连接到NATS服务器: %v", err)
	}
	defer nc.Close()

	// 订阅主题（支持通配符，*匹配一个层级，>匹配所有子层级）
	// 这里订阅所有"app.notifications"相关的消息
	subSubject := "app.notifications"

	// 订阅消息，设置回调函数处理收到的消息
	sub, err := nc.Subscribe(subSubject, func(msg *nats.Msg) {
		// 解析JSON消息
		var pubMsg PubMessage
		if err := json.Unmarshal(msg.Data, &pubMsg); err != nil {
			log.Printf("解析消息失败: %v, 原始数据: %s", err, string(msg.Data))
			return
		}

		// 打印收到的消息详情
		fmt.Printf("\n收到消息: \n")
		fmt.Printf("  主题: %s\n", msg.Subject)
		fmt.Printf("  子主题: %s\n", pubMsg.Topic)
		fmt.Printf("  内容: %s\n", pubMsg.Content)
		fmt.Printf("  时间: %s\n", pubMsg.Time)
	})
	if err != nil {
		log.Fatalf("订阅主题失败: %v", err)
	}
	defer sub.Unsubscribe() // 程序退出时取消订阅

	fmt.Printf("已订阅主题: %s (按Ctrl+C退出)\n", subSubject)
	// 保持程序运行，持续接收消息
	select {}
}
