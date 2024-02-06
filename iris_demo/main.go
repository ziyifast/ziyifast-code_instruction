package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()
	app.Listen(":8081", nil)
}
