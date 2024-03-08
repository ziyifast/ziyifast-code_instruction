package controller

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"myTest/demo_home/blond_filter/service"
	"net/http"
	"strconv"
)

type PlayerController struct {
	Ctx iris.Context
}

func (p *PlayerController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/find/{id}", "FindById")
}

func (p *PlayerController) FindById() mvc.Result {
	defer p.Ctx.Next()
	pId := p.Ctx.Params().Get("id")
	id, err := strconv.ParseInt(pId, 10, 64)
	if err != nil {
		return mvc.Response{
			Code:        http.StatusBadRequest,
			Content:     []byte(err.Error()),
			ContentType: "application/json",
		}
	}
	player, err := service.PlayerService.FindById(id)
	if err != nil {
		return mvc.Response{
			Code:        http.StatusInternalServerError,
			Content:     []byte(err.Error()),
			ContentType: "application/json",
		}
	}
	marshal, err := json.Marshal(player)
	if err != nil {
		return mvc.Response{
			Code:        http.StatusInternalServerError,
			Content:     []byte(err.Error()),
			ContentType: "application/json",
		}
	}
	return mvc.Response{
		Code:        http.StatusOK,
		Content:     marshal,
		ContentType: "application/json",
	}
}
