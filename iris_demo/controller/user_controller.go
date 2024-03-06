package controller

import (
	"github.com/kataras/iris/v12/mvc"
	"github.com/ziyifast/log"
	"net/http"
)

type UserController struct {
	BaseController
}

func (c *UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(http.MethodGet, "/smoke/", "Smoke")
}

func (c *UserController) Smoke() mvc.Result {
	log.Infof("user controller")
	defer c.Ctx.Next()
	return c.Ok("smoke")
}
