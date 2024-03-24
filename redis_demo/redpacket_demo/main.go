package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	context2 "github.com/kataras/iris/v12/context"
	"math/rand"
	"time"
)

/*
通过redis实现迷你版微信抢红包
1. 发红包
2. 拆红包（一个红包拆分成多少个，每个红包里有多少钱）=》二倍均值算法，将拆分后的红包通过list放入redis
3. 抢红包（用户抢红包，并记录哪个用户抢了多少钱，防止重复抢）：hset记录每个红包被哪些用户抢了
*/
var (
	RedisCli                *redis.Client
	RED_PACKGE_KEY          = "redpackage:"
	RED_PACKAGE_CONSUME_KEY = "redpackage:consume:"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	RedisCli = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func main() {
	app := iris.New()
	app.Get("/send", sendRedPacket)
	app.Get("/rob", robRedPacket)
	app.Get("/info", infoRedPacket)
	app.Listen(":9090", nil)
}

// 发红包 http://localhost:9090/send?totalMoney=100&totalNum=3
func sendRedPacket(c *context2.Context) {
	money, _ := c.URLParamInt("totalMoney")
	totalNum, _ := c.URLParamInt("totalNum")
	redPackets := splitRedPacket(money, totalNum)
	uuid, _ := uuid.NewUUID()
	k := RED_PACKGE_KEY + uuid.String()
	for _, r := range redPackets {
		_, err := RedisCli.LPush(context.TODO(), k, r).Result()
		if err != nil && err != redis.Nil {
			panic(err)
		}
	}
	c.JSON(fmt.Sprintf("send redpacket[%s] succ %v", k, redPackets))
}

// 抢红包 http://localhost:9090/rob?redPacket=e3e71f56-e9a3-11ee-9ad5-7a2cb90a4104&uId=4
func robRedPacket(c *context2.Context) {
	//判断是否抢过
	redPacket := c.URLParam("redPacket")
	uId, _ := c.URLParamInt("uId")
	exists, err := RedisCli.HExists(context.TODO(), RED_PACKAGE_CONSUME_KEY+redPacket, fmt.Sprintf("%d", uId)).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if exists {
		//表明已经抢过
		c.JSON(fmt.Sprintf("[%d] you have already rob", uId))
		return
	} else if !exists {
		//从list里取出一个红包
		result, err := RedisCli.LPop(context.TODO(), RED_PACKGE_KEY+redPacket).Result()
		if err == redis.Nil {
			//红包已经抢完了
			c.JSON(fmt.Sprintf("redpacket is empty"))
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d rob the red packet %v\n", uId, result)
		//记录：后续可以异步进MySQL或者MQ做统计分析，每一年抢了多少红包，金额是多少【年度总结】
		_, err = RedisCli.HSet(context.TODO(), RED_PACKAGE_CONSUME_KEY+redPacket, uId, result).Result()
		if err != nil && err != redis.Nil {
			panic(err)
		}
		c.JSON(fmt.Sprintf("[%d] rob the red packet %v", uId, result))
	}

}

func infoRedPacket(c *context2.Context) {
	redPacket := c.URLParam("redPacket")
	infoMap, err := RedisCli.HGetAll(context.TODO(), RED_PACKAGE_CONSUME_KEY+redPacket).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	c.JSON(infoMap)
}

// 拆红包
func splitRedPacket(totalMoney, totalNum int) []int {
	//1. 将红包拆分为几个
	redpackets := make([]int, totalNum)
	usedMoney := 0
	for i := 0; i < totalNum; i++ {
		//最后一个红包，还剩余多少就分多少
		if i == totalNum-1 {
			redpackets[i] = totalMoney - usedMoney
		} else {
			//二倍均值算法：每次拆分后塞进子红包的金额 = 随机区间(0, (剩余红包金额M / 未被抢的剩余红包个数N) * 2)
			avgMoney := ((totalMoney - usedMoney) / (totalNum - i)) * 2
			money := 1 + rand.Intn(avgMoney-1)
			redpackets[i] = money
			usedMoney += money
		}
	}
	return redpackets
}
