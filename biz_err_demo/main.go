package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"myTest/demo_home/biz_err_demo/controller"
)

func main() {
	app := iris.New()
	mvc.New(app).Handle(new(controller.TestBizController))
	app.Listen(":8088", nil)
}
