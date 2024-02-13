package main

import (
	"fmt"
	"sync"
)

// 懒汉式：用到才加载
var (
	instance *Singleton
	once     = sync.Once{}
)

type Singleton struct {
}

func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{}
	})
	return instance
}

func main() {
	one := GetInstance()
	two := GetInstance()
	//one=0x100f54088
	//two=0x100f54088
	fmt.Printf("one=%p\n", one)
	fmt.Printf("two=%p\n", two)
}
