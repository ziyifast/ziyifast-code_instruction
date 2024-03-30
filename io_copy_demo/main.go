package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	app := iris.New()
	go func() {
		http.ListenAndServe(":80", nil)
	}()
	//readAll
	app.Get("/readAll", testReadAll)
	//io.Copy
	app.Get("/ioCopy", func(ctx *context.Context) {
		file, err := os.Open("/Users/ziyi2/GolandProjects/MyTest/demo_home/io_copy_demo/xx.zip")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		_, err = io.Copy(ctx.ResponseWriter(), file)
		if err != nil {
			panic(err)
		}
	})
	app.Listen(":8080", nil)
}

func testReadAll(ctx *context.Context) {
	file, err := os.Open("/Users/ziyi2/GolandProjects/MyTest/demo_home/io_copy_demo/xx.zip")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//simulate onLine err
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	_, err = ctx.Write(bytes)
	if err != nil {
		panic(err)
	}
}
