package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func InitControllers(app *iris.Application) {
	myMvc := mvc.New(app.Party("/user"))
	myMvc.Handle(new(UserController))
}
