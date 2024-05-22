package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/http"
)

func main() {
	app := iris.New()
	app.Use(Cors)
	app.Post("/login", login())
	app.Listen(":8080")
}

func login() func(ctx *context.Context) {
	return func(ctx *context.Context) {
		tmp := make(map[string]interface{})
		err := ctx.ReadJSON(&tmp)
		if err != nil {
			ctx.StatusCode(http.StatusBadRequest)
			return
		}
		username := tmp["username"]
		password := tmp["password"]
		fmt.Println("username:", username, "password:", password)
		if username == "admin" && password == "admin" {
			ctx.StatusCode(http.StatusOK)
			return
		} else {
			ctx.StatusCode(http.StatusForbidden)
			return
		}
	}
}

func Cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	if ctx.Method() == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		ctx.StatusCode(204)
		return
	}
	ctx.Next()
}
