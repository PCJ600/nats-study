package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type PubMessage struct {
	Topic   string `json:"topic"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

func main() {
	// 配置连接参数（含超时和重连）
	opts := []nats.Option{
		nats.Timeout(3 * time.Second),       // 连接超时
		nats.MaxReconnects(2),               // 最大重连次数
	}

	// 建立连接
	nc, err := nats.Connect(nats.DefaultURL, opts...)
	if err != nil {
		log.Fatalf("连接NATS失败: %v", err)
	}
	defer nc.Close()

	subj := "app.notifications"

	for i := 1; i <= 3; i++ {
		// 构建消息
		msg := PubMessage{
			Topic:   fmt.Sprintf("msg-%d", i),
			Content: fmt.Sprintf("第%d条消息", i),
			Time:    time.Now().Format(time.RFC3339),
		}
		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("序列化失败: %v", err)
			continue
		}

		// 调用Publish，直接处理返回的即时错误
		if err := nc.Publish(subj, data); err != nil {
			log.Printf("【即时错误】发布失败 (第%d条): %v", i, err)
			continue // 如连接已关闭，无需继续尝试
		}

		// 强制刷新缓冲区，确保消息被发送（处理网络延迟/故障）
		if err := nc.FlushTimeout(2 * time.Second); err != nil {
			log.Printf("【发送超时】消息未送达服务器 (第%d条): %v", i, err)
			continue
		}

		fmt.Printf("已发布消息: %s\n", string(data))
		time.Sleep(1 * time.Second)
	}

	fmt.Println("发布完成")
}
