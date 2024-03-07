package main

import (
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/ziyifast/log"
	"time"
)

/*
基于redis zset实现延迟队列
*/
var redisdb *redis.Client
var DelayQueueKey = "delay-queue"

func initClient() (err error) {
	redisdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // not set password
		DB:       0,  //use default db
	})
	_, err = redisdb.Ping().Result()
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func main() {
	err := initClient()
	if err != nil {
		log.Errorf("init redis client err: %v", err)
		return
	}
	addTaskToQueue("task1", time.Now().Add(time.Second*3).Unix())
	addTaskToQueue("task2", time.Now().Add(time.Second*8).Unix())
	//执行队列中的任务
	getAndExecuteTask()
}

// executeTime为unix时间戳，作为zset中的score。允许redis按照task应该执行时间来进行排序
func addTaskToQueue(task string, executeTime int64) {
	err := redisdb.ZAdd(DelayQueueKey, redis.Z{
		Score:  float64(executeTime),
		Member: task,
	}).Err()
	if err != nil {
		panic(err)
	}
}

// 从redis中取一个task并执行
func getAndExecuteTask() {
	for {
		tasks, err := redisdb.ZRangeByScore(DelayQueueKey, redis.ZRangeBy{
			Min:    "-inf",
			Max:    fmt.Sprintf("%d", time.Now().Unix()),
			Offset: 0,
			Count:  1,
		}).Result()
		if err != nil {
			time.Sleep(time.Second * 1)
			continue
		}
		//处理任务
		for _, task := range tasks {
			fmt.Println("Execute task: ", task)
			//执行完任务之后用 ZREM 移除该任务
			redisdb.ZRem(DelayQueueKey, task)
		}
		time.Sleep(time.Second * 1)
	}
}
