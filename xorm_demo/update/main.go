package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"myTest/demo_home/xorm_demo/pg"
	"time"
	"xorm.io/xorm"
)

/*
	演示update操作
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
	//1. 根据id更新 engine.ID(user.Id).Update(user)
	//updateById(pg.Engine, &User{Id: 1, Name: "heihei", Age: 53})

	//2. 更新指定字段 engine.ID(user.Id).Cols("age").Update(user)
	updateUserAge(pg.Engine, &User{Id: 1, Age: 23})
}

func updateById(engine *xorm.Engine, user *User) (int64, error) {
	return engine.ID(user.Id).Update(user)
}

func updateUserAge(engine *xorm.Engine, user *User) (int64, error) {
	if user == nil {
		return 0, fmt.Errorf("user cannot be nil")
	}
	return engine.ID(user.Id).Cols("age").Update(user)
}
