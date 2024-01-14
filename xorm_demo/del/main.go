package main

import (
	"github.com/aobco/log"
	_ "github.com/lib/pq"
	"myTest/demo_home/xorm_demo/pg"
	"time"
)

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
	//1. 根据id删除 pg.Engine.Delete(&User{Id: id})
	//deleteById(3)

	//2. 批量删除 pg.Engine.In("id", ids).Delete(new(User))
	deleteByIds([]int64{4, 5})
}

func deleteById(id int64) {
	_, err := pg.Engine.Delete(&User{Id: id})
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	log.Infof("del succ...")
}

func deleteByIds(ids []int64) {
	_, err := pg.Engine.In("id", ids).Delete(new(User))
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	log.Infof("batch del succ...")
}
