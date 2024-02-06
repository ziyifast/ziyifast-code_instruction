package main

import (
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	"myTest/demo_home/swagger_demo/iris/controller"
	_ "myTest/demo_home/swagger_demo/iris/docs"
)

func main() {
	app := iris.New()
	controller.InitControllers(app)
	config := &swagger.Config{
		URL: "http://localhost:8080/swagger/doc.json", //The url pointing to API definition
	}
	app.Get("/swagger/{any}", swagger.CustomWrapHandler(config, swaggerFiles.Handler))
	app.Listen(":8080")
}
