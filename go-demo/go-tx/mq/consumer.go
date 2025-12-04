package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"os"
	"time"
)

func main() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup"), //组，跟服务端对上
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.3.98:9876"})), //地址
	)
	//消费对应topic
	err := c.Subscribe("test", consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range msgs {
				fmt.Printf("订阅回调: %v \n", msgs[i])
			}
			//这个相当于消费者 消息ack，如果失败可以返回 consumer.ConsumeRetryLater
			return consumer.ConsumeSuccess, nil
			//这个相当于失败  要回滚
			//return consumer.ConsumeRetryLater, nil
		})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("关闭消费者失败: %s", err.Error())
	}
}
