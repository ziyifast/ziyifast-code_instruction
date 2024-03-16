package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	context2 "github.com/kataras/iris/v12/context"
	"myTest/demo_home/redis_demo/distributed_lock/constant"
	service2 "myTest/demo_home/redis_demo/distributed_lock/other_svc/service"
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
	//xLock2 := new(sync.Mutex)
	svcCount2 := new(atomic.Int64)
	app.Get("/consume", func(c *context2.Context) {
		//xLock2.Lock()
		//defer xLock2.Unlock()
		service2.GoodsService2.Consume()
		c.JSON("ok port:9999")
	})
	app.Get("/count", func(c *context2.Context) {
		count := svcCount2.Load()
		c.JSON(fmt.Sprintf("%d", count))
	})
	app.Listen(":9999", nil)
}
