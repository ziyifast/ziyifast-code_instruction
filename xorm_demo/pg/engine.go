package pg

import (
	"fmt"
	"github.com/aobco/log"
	"time"
	"xorm.io/xorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbName   = "postgres"
)

//go get "xorm.io/xorm"
var Engine *xorm.Engine

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	engine, err := xorm.NewEngine("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	engine.ShowSQL(true) //菜鸟必备
	engine.SetMaxIdleConns(10)
	engine.SetMaxOpenConns(20)
	engine.SetConnMaxLifetime(time.Minute * 10)
	engine.Cascade(true)
	if err = engine.Ping(); err != nil {
		log.Fatalf("%v", err)
	}
	Engine = engine
	fmt.Println("connect postgresql success")
}
