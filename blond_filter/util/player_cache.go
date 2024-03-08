package util

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/ziyifast/log"
	"myTest/demo_home/blond_filter/model"
	redis2 "myTest/demo_home/blond_filter/redis"
	"strconv"
	"time"
)

type playerCache struct {
}

var (
	PlayerCache = new(playerCache)
	PlayerKey   = "player"
)

func (c *playerCache) GetById(id int64) (*model.Player, error) {
	log.Infof("query redis,time:%v", time.Now().String())
	result, err := redis2.Client.HGet(PlayerKey, strconv.FormatInt(id, 10)).Result()
	if err != nil && err != redis.Nil {
		log.Errorf("%v", err)
		return nil, err
	}
	if result == "" {
		return nil, nil
	}
	p := new(model.Player)
	err = json.Unmarshal([]byte(result), p)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	return p, nil
}

func (c *playerCache) Put(player *model.Player) error {
	marshal, err := json.Marshal(player)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	_, err = redis2.Client.HSet(PlayerKey, strconv.FormatInt(player.Id, 10), string(marshal)).Result()
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	return nil
}
