package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	context2 "github.com/kataras/iris/v12/context"
	"github.com/ziyifast/log"
	"time"
)

var RedisCli *redis.Client

func init() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	_, err := RedisCli.Del(context.TODO(), GoodsKeyCacheA).Result()
	if err != nil {
		panic(err)
	}
}

type Goods struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var (
	GoodsKeyCacheA = "goodsA"
	GoodsKeyCacheB = "goodsB"
)

func InitDb(s, count int, cacheKey string) {
	for i := s; i < s+count; i++ {
		g := &Goods{
			Id:   i + 1,
			Name: fmt.Sprintf("good-%d", i+1),
		}
		marshal, err := json.Marshal(g)
		if err != nil {
			panic(err)
		}
		_, err = RedisCli.RPush(context.TODO(), cacheKey, string(marshal)).Result()
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	InitDb(0, 20, GoodsKeyCacheA)
	app := iris.New()
	app.Get("/goods/top/{offset}/{pageSize}", func(c *context2.Context) {
		offset, err := c.Params().GetInt64("offset")
		if err != nil {
			panic(err)
		}
		pageSize, err := c.Params().GetInt64("pageSize")
		if err != nil {
			panic(err)
		}
		start := (offset - 1) * pageSize
		end := offset*pageSize - 1
		err = c.JSON(QueryForData(start, end))
		if err != nil {
			panic(err)
		}
	})
	//set the expire time
	_, err := RedisCli.Expire(context.TODO(), GoodsKeyCacheA, time.Second*8).Result()
	if err != nil {
		panic(err)
	}
	InitDb(0, 20, GoodsKeyCacheB)
	//add cacheB, expire time is different from cacheA (make sure new goods will be added to cacheA)
	_, err = RedisCli.Expire(context.TODO(), GoodsKeyCacheB, time.Second*20).Result()
	if err != nil {
		panic(err)
	}
	go ReloadNewGoods()
	app.Listen(":9999", nil)
}

func QueryForData(start, end int64) []string {
	val := RedisCli.LRange(context.TODO(), GoodsKeyCacheA, start, end).Val()
	log.Infof("query redis of cache A")
	if len(val) == 0 {
		log.Infof("cacheA is not exist, query redis of cache B")
		val = RedisCli.LRange(context.TODO(), GoodsKeyCacheB, start, end).Val()
		if len(val) == 0 {
			log.Infof("cacheB is not exist, query db, no!!!")
			return val
		}
	}
	return val
}

func ReloadNewGoods() {
	time.Sleep(time.Second * 15)
	log.Infof("start ReloadNewGoods......")
	InitDb(2000, 20, GoodsKeyCacheA)
	//set the expire time of cacheA
	log.Infof("ReloadNewGoods......DONE")
	//reload cacheB....
	//set the expire time of cacheB
}
