package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/ziyifast/log"
	"myTest/demo_home/redis_demo/distributed_lock/constant"
	"myTest/demo_home/redis_demo/distributed_lock/lock"
	"strconv"
)

type goodsService struct {
}

var GoodsService = new(goodsService)

func (g *goodsService) Consume() {
	redisLock := lock.NewRedisLock(constant.RedisCli, constant.BizKey)
	redisLock.Lock()
	defer redisLock.Unlock()
	//consume goods
	result, err := constant.RedisCli.Get(context.TODO(), constant.AppleKey).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	i, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		panic(err)
	}
	if i < 0 {
		log.Infof("no more apple...")
		return
	}
	_, err = constant.RedisCli.Set(context.TODO(), constant.AppleKey, i-1, -1).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	log.Infof("consume success...appleID:%d", i)
}
