package main

import (
	"github.com/aobco/log"
	_ "github.com/lib/pq"
	"myTest/demo_home/xorm_demo/pg"
	"time"
)

/*
	演示insert操作
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
	//1. 插入一条数据 insertOne
	//normalInsert()

	//2. 忽略字段插入 pg.Engine.Omit("user_info").Insert(u)
	//insertWithIgnoreCol()

	//3. 批量插入 pg.Engine.Omit("user_info").Insert(&users)
	insertWithBatch()
}

func normalInsert() {
	u := &User{
		Name: "jack",
		Age:  17,
		UserInfo: Info{
			Address: "beijing",
			Hobbies: []string{
				"baseball",
				"soccer-ball",
			},
		},
	}
	_, err := pg.Engine.InsertOne(u)
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	log.Infof("insert succ..")
}

//2. 忽略字段插入
func insertWithIgnoreCol() {
	u := &User{
		Name: "tom",
		Age:  14,
		UserInfo: Info{
			Address: "sichuan",
			Hobbies: []string{
				"baseball",
				"soccer-ball",
			},
		},
	}
	_, err := pg.Engine.Omit("user_info").Insert(u)
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	log.Infof("insert succ..")
}

//3. 批量插入
func insertWithBatch() {
	u1 := &User{
		Name: "tom1",
		Age:  15,
	}
	u2 := &User{
		Name: "tom2",
		Age:  16,
	}
	u3 := &User{
		Name: "tom3",
		Age:  17,
	}
	users := []*User{u1, u2, u3}
	_, err := pg.Engine.Omit("user_info").Insert(&users)
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	log.Infof("batch insert succ..")
}
