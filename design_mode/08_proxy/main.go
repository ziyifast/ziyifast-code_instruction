package main

import (
	"fmt"
)

type PaymentService interface {
	pay(order string) string
}

// AliPay 阿里支付类
type AliPay struct {
}

/**
 * @Description: 阿里支付类，从阿里获取支付token
 * @receiver a
 * @param order
 * @return string
 */
func (a *AliPay) pay(order string) string {
	return "从阿里获取支付token"
}

type PaymentProxy struct {
	realPay PaymentService
}

/**
 * @Description: 做校验签名、初始化订单数据、参数检查、记录日志、组装这种通用性操作，调用真正支付类获取token
 * @receiver p
 * @param order
 * @return string
 */
func (p *PaymentProxy) pay(order string) string {
	fmt.Println("处理" + order)
	fmt.Println("1校验签名")
	fmt.Println("2格式化订单数据")
	fmt.Println("3参数检查")
	fmt.Println("4记录请求日志")
	token := p.realPay.pay(order)
	return "http://组装" + token + "然后跳转到第三方支付"
}
func main() {
	proxy := &PaymentProxy{
		realPay: &AliPay{},
	}
	url := proxy.pay("阿里订单")
	fmt.Println(url)
}
