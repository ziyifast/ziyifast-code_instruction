package controller

import (
	"errors"
	"github.com/kataras/iris/v12/mvc"
	"myTest/demo_home/biz_err_demo/error/biz_err"
	"myTest/demo_home/biz_err_demo/error/zerr"
	"myTest/demo_home/biz_err_demo/response"
	"net/http"
)

type TestBizController struct {
	BaseController
}

func (t *TestBizController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(http.MethodGet, "/testBizErr", "TestBizErr")
}

func (t *TestBizController) TestBizErr() mvc.Result {
	err1 := errors.New("")
	err := zerr.BizWrap(err1, biz_err.UsernameOrPasswordInValid, "")
	return response.JsonBizError(err)
}
