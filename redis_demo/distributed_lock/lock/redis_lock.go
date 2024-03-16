package lock

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/ziyifast/log"
	"strings"
	"time"
)

/*
	通过Redis+Lua脚本实现分布式锁
	一个分布式锁具备的条件：
	1.排他性
	2.高可用
	3.防止死锁（过期时间）、自动续期
	4.可重入
	5.不乱抢，只能删除自己的锁
*/

// lua脚本用于保证Redis加锁&设置过期时间的原子性操作
//EX second ：设置键的过期时间为 second 秒。 SET key value EX second 效果等同于 SETEX key second value 。
//PX millisecond ：设置键的过期时间为 millisecond 毫秒。 SET key value PX millisecond 效果等同于 PSETEX key millisecond value 。
//NX ：只在键不存在时，才对键进行设置操作。 SET key value NX 效果等同于 SETNX key value 。
//XX ：只在键已经存在时，才对键进行设置操作。

var (
	defaultExpireTime = 5 //单位：s
)

type RedisLock struct {
	key string
	// 锁的过期时间，单位: s
	expire uint32
	// 锁的标识
	Id string
	// Redis客户端
	redisCli *redis.Client
}

func NewRedisLock(cli *redis.Client, key string) *RedisLock {
	//去掉uuid中间的-
	id := strings.Join(strings.Split(uuid.New().String(), "-"), "")
	return &RedisLock{
		key:      key,
		expire:   uint32(defaultExpireTime),
		Id:       id,
		redisCli: cli,
	}
}

func (r *RedisLock) TryLock() bool {
	//通过lua脚本加锁[hincrby如果key不存在，则会主动创建,如果存在则会给count数加1，表示又重入一次]
	lockCmd := "if redis.call('exists', KEYS[1]) == 0 or redis.call('hexists', KEYS[1], ARGV[1]) == 1 " +
		"then " +
		"   redis.call('hincrby', KEYS[1], ARGV[1], 1) " +
		"   redis.call('expire', KEYS[1], ARGV[2]) " +
		"   return 1 " +
		"else " +
		"   return 0 " +
		"end"
	result, err := r.redisCli.Eval(context.TODO(), lockCmd, []string{r.key}, r.Id, r.expire).Result()
	if err != nil {
		log.Errorf("tryLock %s %v", r.key, err)
		return false
	}
	i := result.(int64)
	if i == 1 {
		//获取锁成功&自动续期
		go r.reNewExpire()
		return true
	}
	return false
}

func (r *RedisLock) Lock() {
	for {
		if r.TryLock() {
			break
		}
		time.Sleep(time.Millisecond * 20)
	}
}

func (r *RedisLock) SetExpire(t uint32) {
	r.expire = t
}

func (r *RedisLock) Unlock() {
	//通过lua脚本删除锁
	//1. 查看锁是否存在，如果不存在，直接返回
	//2. 如果存在，对锁进行hincrby -1操作,当减到0时，表明已经unlock完成，可以删除key
	delCmd := "if redis.call('hexists', KEYS[1], ARGV[1]) == 0 " +
		"then " +
		"   return nil " +
		"elseif redis.call('hincrby', KEYS[1], ARGV[1], -1) == 0 " +
		"then " +
		"   return redis.call('del', KEYS[1]) " +
		"else " +
		"   return 0 " +
		"end"
	resp, err := r.redisCli.Eval(context.TODO(), delCmd, []string{r.key}, r.Id).Result()
	if err != nil && err != redis.Nil {
		log.Errorf("unlock %s %v", r.key, err)
	}
	if resp == nil {
		fmt.Println("delKey=", resp)
		return
	}
}

// 自动续期
func (r *RedisLock) reNewExpire() {
	renewCmd := "if redis.call('hexists', KEYS[1], ARGV[1]) == 1 " +
		"then " +
		"   return redis.call('expire', KEYS[1], ARGV[2]) " +
		"else " +
		"   return 0 " +
		"end"
	ticker := time.NewTicker(time.Duration(r.expire/3) * time.Second)
	for {
		select {
		case <-ticker.C:
			//查看锁是否存在，如果存在进行续期
			resp, err := r.redisCli.Eval(context.TODO(), renewCmd, []string{r.key}, r.Id, r.expire).Result()
			if err != nil && err != redis.Nil {
				log.Errorf("renew key %s err %v", r.key, err)
			}
			if resp.(int64) == 0 {
				return
			}
			log.Infof("renew.....ing...")
		}
	}
}
