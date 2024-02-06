package controller

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"myTest/demo_home/swagger_demo/iris/model"
	"net/http"
)

type UserController struct {
	Ctx iris.Context
}

func (u *UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(http.MethodGet, "/getAll", "GetAllUsers")
}

// GetAllUsers @Summary 获取用户信息
// @Description 获取所有用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Router /user/getAll [get]
func (u *UserController) GetAllUsers() mvc.Result {
	//手动模拟从数据库查询到user信息
	resp := new(mvc.Response)
	resp.ContentType = "application/json"
	user1 := new(model.User)
	user1.Name = "zhangsan"
	user1.Age = 20
	user2 := new(model.User)
	user2.Name = "li4"
	user2.Age = 28
	users := []model.User{*user1, *user2}
	marshal, _ := json.Marshal(users)
	resp.Content = marshal
	resp.Code = http.StatusOK
	return resp
}
