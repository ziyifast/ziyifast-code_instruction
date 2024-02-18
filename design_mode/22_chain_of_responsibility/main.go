package main

import "fmt"

/*
	抽象处理器（Handler）：定义出一个处理请求的接口。
	具体处理器（ConcreteHandler）：实现抽象处理器的接口，处理它所负责的请求。如果不处理该请求，则把请求转发给它的后继者。
	客户端（Client）：创建处理器对象，并将请求发送到某个处理器。

	案例：我们在处理一些法律案件时，是层层上报，如果乡镇处理不了就交给市区处理，如果市区处理不了就交给省处理，如果省处理不了就交给国家处理...
*/

type Handler interface {
	SetNext(handler Handler) //设置下一个处理器
	Handle(request int)      //处理请求
}

type TownHandler struct {
	NextHandler Handler
}

func (t *TownHandler) SetNext(handler Handler) {
	t.NextHandler = handler
}

func (t *TownHandler) Handle(request int) {
	//处理刑事案件，案件级别小于20
	if request <= 20 {
		fmt.Println("TownHandler: 小于等于20，我来处理。")
	} else {
		if t.NextHandler != nil {
			t.NextHandler.Handle(request)
		}
	}
}

type CityHandler struct {
	NextHandler Handler
}

func (c *CityHandler) SetNext(handler Handler) {
	c.NextHandler = handler
}

func (c *CityHandler) Handle(request int) {
	if request > 20 && request <= 100 {
		fmt.Println("CityHandler: 大于20小于等于100，我来处理。")
	} else {
		if c.NextHandler != nil {
			c.NextHandler.Handle(request)
		}
	}
}

type ProvinceHandler struct {
	NextHandler Handler
}

func (p *ProvinceHandler) SetNext(handler Handler) {
	p.NextHandler = handler
}

func (p *ProvinceHandler) Handle(request int) {
	if request > 100 {
		fmt.Println("ProvinceHandler: 大于100，我来处理。")
	} else {
		if p.NextHandler != nil {
			p.NextHandler.Handle(request)
		}
	}
}

func main() {
	townHandler := &TownHandler{}
	cityHandler := &CityHandler{}
	provinceHandler := &ProvinceHandler{}

	townHandler.SetNext(cityHandler)
	cityHandler.SetNext(provinceHandler)

	// 处理请求
	requests := []int{5, 50, 300}
	for _, request := range requests {
		townHandler.Handle(request)
	}
}
