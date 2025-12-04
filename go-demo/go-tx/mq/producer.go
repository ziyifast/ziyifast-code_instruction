package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// 事务消息的结构体
type DemoListener struct {
	localTrans       *sync.Map
	transactionIndex int32
}

func NewDemoListener() *DemoListener {
	return &DemoListener{
		localTrans: new(sync.Map),
	}
}

// 执行并发送
func (dl *DemoListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	nextIndex := atomic.AddInt32(&dl.transactionIndex, 1)
	fmt.Printf("nextIndex: %v for transactionID: %v\n", nextIndex, msg.TransactionId)
	status := nextIndex % 3
	dl.localTrans.Store(msg.TransactionId, primitive.LocalTransactionState(status+1))
	//在执行SendMessageInTransaction方法的时候会调用此方法ExecuteLocalTransaction，
	//如果ExecuteLocalTransaction 返回primitive.UnknowState 那么brocker就会调用CheckLocalTransaction方法检查消息状态
	// 如果返回  primitive.CommitMessageState 和primitive.RollbackMessageState 则不会调用CheckLocalTransaction
	return primitive.UnknowState
	//return primitive.RollbackMessageState
	//return primitive.CommitMessageState
}

// 检查本地事务是否成功
func (dl *DemoListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	v, existed := dl.localTrans.Load(msg.TransactionId)
	if !existed {
		fmt.Printf("unknow msg: %v, return Commit", msg)
		return primitive.CommitMessageState
	}
	state := v.(primitive.LocalTransactionState)
	fmt.Printf("检查本地事务是否成功 msg transactionID : %v\n", msg.TransactionId)
	switch state {
	case 1:
		fmt.Printf("回滚: %v\n", msg.Body)
		return primitive.RollbackMessageState
	case 2:
		fmt.Printf("未知: %v\n", msg.Body)
		return primitive.UnknowState
	default:
		fmt.Printf("默认提交: %v\n", msg.Body)
		return primitive.CommitMessageState
	}
}

func main() {
	p, _ := rocketmq.NewTransactionProducer(
		NewDemoListener(), //自定义listener
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.3.98:9876"})), //连接地址
		producer.WithRetry(1),               //重试次数
		producer.WithGroupName("testGroup"), //生产者组，可有可无，跟client要对上
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("开启失败: %s\n", err.Error())
		os.Exit(1)
	}

	topic := "test"
	for i := 0; i < 10; i++ {
		res, err := p.SendMessageInTransaction(
			context.Background(),
			primitive.NewMessage(topic, []byte("测试RocketMQ事务消息"+strconv.Itoa(i))),
		)

		if err != nil {
			fmt.Printf("发送消息失败: %s\n", err)
		} else {
			fmt.Printf("发送消息成功=%s\n", res.String())
		}
	}
	time.Sleep(5 * time.Minute)
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("关闭失败: %s", err.Error())
	}
}
