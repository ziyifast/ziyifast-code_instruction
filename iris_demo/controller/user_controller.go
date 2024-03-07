package controller

import (
	"github.com/kataras/iris/v12/mvc"
	"github.com/ziyifast/log"
	"myTest/demo_home/iris_demo/constant"
	"net/http"
)

type UserController struct {
	BaseController
}

func (c *UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(http.MethodGet, "/smoke", "Smoke")
	b.Handle(http.MethodGet, "/login", "Login")
	b.Handle(http.MethodGet, "/logout", "Logout")
}

func (c *UserController) Smoke() mvc.Result {
	session := constant.SessionMgr.Start(c.Ctx)
	log.Infof("sessionID:%v", session.ID())
	defer c.Ctx.Next()
	log.Infof("%v", session.Get("authenticated"))
	if session.Get("authenticated") == nil {
		return c.Unauthorized()
	}
	log.Infof("user controller")
	return c.Ok("smoke")
}

func (c *UserController) Login() mvc.Result {
	//verify the username and password
	session := constant.SessionMgr.Start(c.Ctx)
	session.Set("authenticated", "true")
	return c.Ok("login success")
}

func (c *UserController) Logout() mvc.Result {
	//verify the username and password
	session := constant.SessionMgr.Start(c.Ctx)
	session.Set("authenticated", nil)
	return c.Ok("login out...")
}
