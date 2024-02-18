package main

import "fmt"

/*
	通过访问者模式实现在生产环境和开发环境中打印处不同的内容
*/

type IVisitor interface {
	Visit() //访问者的访问方法
}

type ProductionVisitor struct {
	env string
}

func (p *ProductionVisitor) Visit() {
	if p.env == "pro" {
		fmt.Println("这是生产环境的输出")
	}
}

type DevelopmentVisitor struct {
	env string
}

func (d *DevelopmentVisitor) Visit() {
	if d.env == "dev" {
		fmt.Println("这是开发环境的输出")
	}
}

// IElement IElement接口，在其中声明一个accept()操作，它以一个抽象访问者作为参数
type IElement interface {
	Accept(visitor IVisitor)
}

// Element 具体元素，它实现了accept()操作，在accept()中调用访问者的访问方法以便完成对一个元素的操作
type Element struct {
	visitors []IVisitor
}

func (p *Element) Accept(visitor IVisitor) {
	p.visitors = append(p.visitors, visitor)
}

type ExampleLog struct {
	Element
}

func (e *ExampleLog) Print() {
	for _, visitor := range e.visitors {
		visitor.Visit()
	}
}

func main() {
	ele := Element{}
	env := "dev"
	ele.Accept(&ProductionVisitor{env: env})
	ele.Accept(&DevelopmentVisitor{env: env})

	example := &ExampleLog{Element: ele}
	example.Print()
}
