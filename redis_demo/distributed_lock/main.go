package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	context2 "github.com/kataras/iris/v12/context"
	"myTest/demo_home/redis_demo/distributed_lock/constant"
	"myTest/demo_home/redis_demo/distributed_lock/service"
	"sync/atomic"
)

func main() {
	constant.RedisCli = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	_, err := constant.RedisCli.Set(context.TODO(), constant.AppleKey, 500, -1).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	app := iris.New()
	//xLock := new(sync.Mutex)
	svcCount := new(atomic.Int64)
	app.Get("/consume", func(c *context2.Context) {
		svcCount.Add(1)
		//xLock.Lock()
		//defer xLock.Unlock()
		service.GoodsService.Consume()

		c.JSON("ok port:8888")
	})
	app.Get("/count", func(c *context2.Context) {
		fmt.Println("middle...")
		count := svcCount.Load()
		c.JSON(fmt.Sprintf("%d", count))
	})
	app.Listen(":8888", nil)
}
