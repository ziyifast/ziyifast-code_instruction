package main

import (
	"github.com/aobco/log"
	_ "github.com/lib/pq"
	"myTest/demo_home/xorm_demo/pg"
	"time"
)

/*
	通过xorm自动创建表结构
*/

type User struct {
	Id          int64     `xorm:"bigint pk autoincr"`
	Name        string    `xorm:"varchar(25) notnull unique comment('姓名')"`
	Age         int64     `xorm:"bigint"`
	UserInfo    Info      `xorm:"JSON"`
	CreatedTime time.Time `xorm:"created_time timestampz created"`
	ModifyTime  time.Time `xorm:"modify_time timestampz updated"`
}

type Info struct {
	Address string   `json:"address"`
	Hobbies []string `json:"hobbies"`
}

/*
 xorm官网文档：https://xorm.io/zh/docs
*/

func main() {
	err := pg.Engine.Sync(new(User))
	if err != nil {
		log.Errorf("同步表结构失败")
		return
	}
	log.Infof("同步表结构成功...")
}
