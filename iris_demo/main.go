package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/ziyifast/log"
	"myTest/demo_home/iris_demo/constant"
	"myTest/demo_home/iris_demo/controller"
	"myTest/demo_home/iris_demo/route"
	"myTest/demo_home/iris_demo/util"
	"time"
)

var (
	signaturePrefix = "/yi/sign/"
)

func main() {
	InitControllers()
	app := iris.New()
	initMvcHandle(app)
	app.Listen(":8899")
}
func initSession() {
	sessionCfg := sessions.Config{
		Cookie:  "test",
		Expires: time.Duration(10) * time.Second,
	}
	constant.SessionMgr = sessions.New(sessionCfg)
}

func initMvcHandle(app *iris.Application) {
	initSession()
	for _, v := range route.ControllerList {
		log.Debugf("routeName:%s middleware:%v  doneHandler:%v", v.RouteName, v.MiddlewareSlice, v.DoneHandleSlice)
		myMvc := mvc.New(app.Party(v.RouteName))
		myMvc.Router.Use(v.MiddlewareSlice...)
		myMvc.Router.Done(v.DoneHandleSlice...)
		myMvc.Register(constant.SessionMgr.Start)
		myMvc.Handle(v.ControllerObj)
	}
}

func InitControllers() {
	util.NewSignRoute(signaturePrefix, "user", new(controller.UserController))
}
