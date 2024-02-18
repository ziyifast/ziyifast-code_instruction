package main

import "fmt"

// ActionState 定义状态接口：每个状态可以对应那些动作
type ActionState interface {
	View()
	Comment()
	Post()
}

// Account 定义账户结构体：包含当前账户状态State、HealthValue账号健康值
type Account struct {
	State       ActionState
	HealthValue int
}

func NewAccount(health int) *Account {
	a := new(Account)
	a.SetHealth(health)
	a.changeState()
	return a
}

func (a *Account) Post() {
	a.State.Post()
}

func (a *Account) View() {
	a.State.View()
}

func (a *Account) Comment() {
	a.State.Comment()
}

type NormalState struct {
}

func (n *NormalState) Post() {
	fmt.Println("正常发帖")
}

func (n *NormalState) View() {
	fmt.Println("正常看帖")
}

func (n *NormalState) Comment() {
	fmt.Println("正常评论")
}

type RestrictState struct {
}

func (r *RestrictState) Post() {
	fmt.Println("抱歉，你的健康值小于0，不能发帖")
}

func (r *RestrictState) View() {
	fmt.Println("正常看帖")
}

func (r *RestrictState) Comment() {
	fmt.Println("正常评论")
}

type CloseState struct {
}

func (c *CloseState) Post() {
	fmt.Println("抱歉，你的健康值小于0，不能发帖")
}

func (c *CloseState) View() {
	fmt.Println("账号被封，无法看帖")
}

func (c *CloseState) Comment() {
	fmt.Println("抱歉，你的健康值小于-10，不能评论")
}

func (a *Account) SetHealth(value int) {
	a.HealthValue = value
	a.changeState()
}

func (a *Account) changeState() {
	if a.HealthValue <= -10 {
		a.State = &CloseState{}
	} else if a.HealthValue > -10 && a.HealthValue <= 0 {
		a.State = &RestrictState{}
	} else if a.HealthValue > 0 {
		a.State = &NormalState{}
	}
}

func main() {
	fmt.Println("===========正常账户===========")
	//正常账户：可发帖、评论、查看
	account := NewAccount(10)
	account.Post()
	account.View()
	account.Comment()
	fmt.Println("===========受限账户===========")
	//受限账户：不能发帖、可以评论和查看
	account.SetHealth(-5)
	account.Post()
	account.View()
	account.Comment()
	fmt.Println("===========被封号账户===========")
	//被封号账户：不能发帖、不能评论、不能查看
	account.SetHealth(-11)
	account.Post()
	account.View()
	account.Comment()
}
