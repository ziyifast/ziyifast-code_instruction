package main

import "fmt"

// AlipayInterface 支付宝SDK
type AlipayInterface interface {
	Pay(money int)
}

type AlipayPay struct {
}

func (a *AlipayPay) Pay(money int) {
	fmt.Println("这里是支付宝支付：", "费用是：", money)
}

type WeixinPayInterface interface {
	Pay(money int)
}
type WeixinPay struct {
}

func (a *WeixinPay) Pay(money int) {
	fmt.Println("这里是微信支付：", "费用是：", money)
}

// TargetPayInterface 目标接口，能支持传入支付宝或者微信支付进行支付
type TargetPayInterface interface {
	DealPay(payType string, money int)
}

// 自己的adapter，实现微信和支付宝支付，
type NewAdapter struct {
	AlipayInterface
	WeixinPayInterface
}

func (n *NewAdapter) DealPay(payType string, money int) {
	switch payType {
	case "weixinpay":
		n.WeixinPayInterface.Pay(money)
	case "alipay":
		n.AlipayInterface.Pay(money)
	default:
		fmt.Println("不支持的支付方式")
	}
}

func main() {
	// 同时调用支付宝和微信支付
	t := &NewAdapter{
		AlipayInterface:    &AlipayPay{},
		WeixinPayInterface: &WeixinPay{},
	}
	// 这里业务中基于一个用户同时只能调用一种支付方式。
	t.DealPay("weixinpay", 35)
	t.DealPay("alipay", 101)
}
