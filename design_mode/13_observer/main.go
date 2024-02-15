package main

import "fmt"

/*
	被观察者持有了集合存放观察者 (收通知的为观察者)
	- 报纸订阅，报社为被观察者，订阅的人为观察者
	- MVC 模式，当 model 改变时，View 视图会自动改变，model 为被观察者，View 为观察者
*/

type Customer interface {
	update()
}

type CustomerA struct {
}

func (c *CustomerA) update() {
	fmt.Println("我是客户A, 我收到报纸了")
}

type CustomerB struct {
}

func (c *CustomerB) update() {
	fmt.Println("我是客户B, 我收到报纸了")
}

// NewsOffice 被观察者（报社）：存储了所有的观察者，订阅报纸的人为观察者
type NewsOffice struct {
	customers []Customer
}

func (n *NewsOffice) addCustomer(c Customer) {
	n.customers = append(n.customers, c)
}

func (n *NewsOffice) newspaperCome() {
	//新报纸到了，通知所有客户（观察者）
	fmt.Println("新报纸到了，通知所有客户（观察者）")
	n.notifyAllCustomer()
}

func (n *NewsOffice) notifyAllCustomer() {
	for _, customer := range n.customers {
		customer.update()
	}
}

func main() {
	a := &CustomerA{}
	b := &CustomerB{}
	office := NewsOffice{}
	// 模拟客户订阅
	office.addCustomer(a)
	office.addCustomer(b)
	// 新的报纸
	office.newspaperCome()
}
