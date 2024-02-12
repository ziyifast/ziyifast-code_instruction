package main

import "fmt"

// Prototype 定于原型接口
type Prototype interface {
	Clone() Prototype
}

// ConcretePrototype 定义具体原型结构体
type ConcretePrototype struct {
	name string
}

// Clone 提供clone方法
func (c ConcretePrototype) Clone() Prototype {
	return &ConcretePrototype{
		name: c.name,
	}
}

func main() {
	//创建原型对象
	prototypeObj := &ConcretePrototype{name: "test"}
	//使用原型对象创建新对象
	cloneObject := prototypeObj.Clone()
	fmt.Println(cloneObject.(*ConcretePrototype).name)
	//prototypeObj=0x1400000e028
	//cloneObject=0x14000010280
	fmt.Printf("prototypeObj=%v\n", &prototypeObj)
	fmt.Printf("cloneObject=%p\n", cloneObject)
}
