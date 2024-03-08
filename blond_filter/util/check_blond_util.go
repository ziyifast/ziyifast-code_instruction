package util

import (
	"fmt"
	"github.com/ziyifast/log"
	"math"
	"myTest/demo_home/blond_filter/redis"
)

var base = 1 << 32

// achieve blond filter
// 1. calculate the hash of key
// 2. preload the players data
func InitBlondFilter() {
	//get hashCode
	key := fmt.Sprintf("%s%d", redis.PlayerPrefix, 1)
	hashCode := int(math.Abs(float64(getHashCode(key))))
	//calculate the offset
	offset := hashCode % base
	_, err := redis.Client.SetBit(key, int64(offset), 1).Result()
	if err != nil {
		panic(err)
	}
}

func getHashCode(str string) int {
	var hash int32 = 17
	for i := 0; i < len(str); i++ {
		hash = hash*31 + int32(str[i])
	}
	return int(hash)
}

func CheckExist(id int64) bool {
	key := fmt.Sprintf("%s%d", redis.PlayerPrefix, id)
	hashCode := int(math.Abs(float64(getHashCode(key))))
	offset := hashCode % base
	res, err := redis.Client.GetBit(key, int64(offset)).Result()
	if err != nil {
		log.Errorf("%v", err)
		return false
	}
	log.Infof("%v", res)
	return res == 1
}
