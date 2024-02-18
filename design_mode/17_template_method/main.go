package main

import "fmt"

/*
	关键：通用步骤在抽象类中实现，变化的步骤在具体的子类中实现
	例：做饭，打开煤气，开火，（做饭）， 关火，关闭煤气。除了做饭其他步骤都是相同的，抽到抽象类中实现
	1. 定义接口，包含做饭的全部步骤
	2. 定义抽象类type CookMenu struct，实现做饭的通用步骤，打开煤气、开火...
	3. 定义具体子类XiHongShi、ChaoJiDan struct，重写cook方法（不同的子类有不同的实现）
*/

type Cooker interface {
	open()
	fire()
	cook()
	closefire()
	close()
}

// 类似于一个抽象类
type CookMenu struct {
}

func (CookMenu) open() {
	fmt.Println("打开开关")
}

func (CookMenu) fire() {
	fmt.Println("开火")
}

// 做菜，交给具体的子类实现
func (CookMenu) cooke() {
}

func (CookMenu) closefire() {
	fmt.Println("关火")
}

func (CookMenu) close() {
	fmt.Println("关闭开关")
}

// 封装具体步骤
func doCook(cook Cooker) {
	cook.open()
	cook.fire()
	cook.cook()
	cook.closefire()
	cook.close()
}

type XiHongShi struct {
	CookMenu
}

func (*XiHongShi) cook() {
	fmt.Println("做西红柿")
}

type ChaoJiDan struct {
	CookMenu
}

func (ChaoJiDan) cook() {
	fmt.Println("做炒鸡蛋")
}

func main() {
	x := &XiHongShi{}
	doCook(x)
	fmt.Println("============")
	y := &ChaoJiDan{}
	doCook(y)
}
