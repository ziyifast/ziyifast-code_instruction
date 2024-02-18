package main

import "fmt"

// PaymentStrategy 定义策略接口
type PaymentStrategy interface {
	Pay(amount float64) error
}

// CreditCardStrategy 实现具体的支付策略：信用卡支付
type CreditCardStrategy struct {
	name     string
	cardNum  string
	password string
}

func (c *CreditCardStrategy) Pay(amount float64) error {
	fmt.Printf("Paying %0.2f using credit card\n", amount)
	return nil
}

// CashStrategy 实现具体的支付策略：现金支付
type CashStrategy struct {
	name string
}

func (c *CashStrategy) Pay(amount float64) error {
	fmt.Printf("Paying %0.2f by cash \n", amount)
	return nil
}

// PaymentContext 定义上下文类
type PaymentContext struct {
	amount   float64
	strategy PaymentStrategy
}

// Pay 封装pay方法：通过调用strategy的pay方法
func (p *PaymentContext) Pay() error {
	return p.strategy.Pay(p.amount)
}

func NewPaymentContext(amount float64, strategy PaymentStrategy) *PaymentContext {
	return &PaymentContext{
		amount:   amount,
		strategy: strategy,
	}
}

func main() {
	creditCardStrategy := &CreditCardStrategy{
		name:    "John Doe",
		cardNum: "1234 5678 9012 3456",
	}
	paymentContext := NewPaymentContext(20.0, creditCardStrategy)
	paymentContext.Pay()
	cashStrategy := &CashStrategy{
		name: "Juicy",
	}
	cashPaymentContext := NewPaymentContext(110.0, cashStrategy)
	cashPaymentContext.Pay()
}
