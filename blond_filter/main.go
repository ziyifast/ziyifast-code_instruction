package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"myTest/demo_home/blond_filter/controller"
	"myTest/demo_home/blond_filter/util"
)

func main() {
	//pg.Engine.Sync(new(model.Player))
	app := iris.New()
	pMvc := mvc.New(app.Party("player"))
	pMvc.Handle(new(controller.PlayerController))
	util.InitBlondFilter()
	app.Listen(":9999", nil)
}
