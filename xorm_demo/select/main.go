package main

import (
	"errors"
	"github.com/aobco/log"
	_ "github.com/lib/pq"
	"myTest/demo_home/xorm_demo/pg"
	"time"
	"xorm.io/xorm"
)

/*
	演示查询
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
	//user, _ := getById(pg.Engine, 1)
	//log.Infof("%v", user)

	//users, _ := getByIds(pg.Engine, []int64{1, 2})
	//for _, user := range users {
	//	log.Infof("%+v", user)
	//}
	//
	//users, _ := getByName(pg.Engine, "tom")
	//for _, user := range users {
	//	log.Infof("%+v", user)
	//}

	users, _ := listByAge(pg.Engine, 10)
	for _, user := range users {
		log.Infof("%+v", user)
	}
}

//1. 根据id查询
func getById(engine xorm.Interface, id int64) (*User, error) {
	u := new(User)
	flag, err := engine.ID(id).Get(u)
	if err != nil {
		return nil, err
	}
	if !flag {
		return nil, nil
	}
	return u, nil
}

//2. 范围查询
func getByIds(engine *xorm.Engine, ids []int64) ([]*User, error) {
	users := make([]*User, 0)
	if len(ids) == 0 {
		return nil, errors.New("ids is required")
	}
	err := engine.In("id", ids).Find(&users)
	return users, err
}

//3. 模糊查询[单条记录用Get、多条记录用Find]
func getByName(engine *xorm.Engine, name string) ([]*User, error) {
	users := make([]*User, 0)
	err := engine.Where("name like ? ", "%"+name+"%").Find(&users)
	return users, err
}

func listByAge(engine *xorm.Engine, age int64) ([]*User, error) {
	users := make([]*User, 0)
	//err := engine.Where("age > ?", age).Decr("created_time").Find(&users)
	//err := engine.Where("age > ?", age).OrderBy("created_time").Find(&users)
	err := engine.Where("age > ?", age).OrderBy("created_time").Limit(2, 0).Find(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
