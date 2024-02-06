package main

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"myTest/demo_home/swagger_demo/gin/controller"
	_ "myTest/demo_home/swagger_demo/gin/docs"
)

var swagHandler gin.HandlerFunc

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server.
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// SwaggerUI: http://localhost:8080/swagger/index.html
func main() {
	e := gin.Default()
	v1 := e.Group("/api/v1")
	{
		v1.GET("/hello", controller.Hello)
		v1.POST("/login", controller.Login)
	}
	e.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if swagHandler != nil {
		e.GET("/swagger/*any", swagHandler)
	}

	if err := e.Run(":8080"); err != nil {
		panic(err)
	}
}
