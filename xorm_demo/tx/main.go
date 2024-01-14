package main

import (
	"errors"
	"github.com/aobco/log"
	_ "github.com/lib/pq"
	"myTest/demo_home/xorm_demo/pg"
	"time"
	"xorm.io/xorm"
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

/*
	xorm操作事物
*/
func main() {
	session := pg.Engine.NewSession()
	session.Begin()
	defer func() {
		session.Commit()
	}()
	u := new(User)
	u.Id = 7
	u.Age = -18
	_, err := session.ID(u.Id).Cols("age").Update(u)
	if err != nil {
		log.Errorf("Failed to update user 1: %v", err)
		session.Rollback()
		return
	}
	// 模拟构造一个错误
	err = errors.New("built an error")
	if err != nil {
		log.Errorf("Simulated error: %v", err)
		session.Rollback()
		return
	}

	u2 := new(User)
	u2.Id = 9
	u2.Age = 555
	_, err = session.ID(u2.Id).Cols("age").Update(u2)
	if err != nil {
		log.Errorf("Failed to update user 2: %v", err)
		session.Rollback()
		return
	}
	return
}

func updateUserInfo(session *xorm.Session) error {
	u := new(User)
	u.Id = 7
	u.Age = 9
	_, err := session.ID(u.Id).Cols("age").Update(u)
	//err = errors.New("build an error")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	u2 := new(User)
	u2.Id = 9
	u2.Age = 22222
	_, err = session.ID(u2.Id).Cols("age").Update(u2)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	return err
}

type SessionHandleFunc func(session *xorm.Session) error

func Tx(f SessionHandleFunc) error {
	session := pg.Engine.NewSession()
	session.Begin()
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%+v", err)
			session.Rollback()
		}
	}()
	err := f(session)
	if err != nil {
		log.Errorf("[DB_TX] error %+v", err)
		session.Rollback()
		return err
	}
	return session.Commit()
}
