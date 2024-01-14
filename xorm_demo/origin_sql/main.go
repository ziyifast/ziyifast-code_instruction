package main

import (
	"github.com/aobco/log"
	_ "github.com/lib/pq"
	"myTest/demo_home/xorm_demo/pg"
	"time"
)

/*
	演示原生SQL操作
*/

type User struct {
	Id          int64     `xorm:"bigint pk autoincr"`
	Name        string    `xorm:"varchar(25) notnull unique comment('姓名')"`
	Age         int64     `xorm:"bigint"`
	UserInfo    Info      `xorm:"user_info JSON"`
	CreatedTime time.Time `xorm:"created_time timestampz created"`
	ModifyTime  time.Time `xorm:"modify_time timestampz updated"`
}

type Info struct {
	Address string   `json:"address"`
	Hobbies []string `json:"hobbies"`
}

func main() {
	name := "jack"
	u := new(User)
	_, err := pg.Engine.SQL("select id,name,age from public.user where name = ?", name).Get(u)
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	log.Infof("res: %+v", u)
}
