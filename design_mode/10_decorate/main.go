package main

import "fmt"

type Car struct {
	Brand string
	Price float64
}

// PriceDecorator 定义装饰器接口
type PriceDecorator interface {
	DecoratePrice(c Car) Car
}

// ExtraPriceDecorator 实现装饰器
type ExtraPriceDecorator struct {
	ExtraPrice float64
}

func (d ExtraPriceDecorator) DecoratePrice(car Car) Car {
	car.Price += d.ExtraPrice
	return car
}

func main() {
	toyota := Car{Brand: "Toyota", Price: 10000}
	decorator := ExtraPriceDecorator{ExtraPrice: 500}
	decoratedCar := decorator.DecoratePrice(toyota)
	fmt.Printf("%+v\n", decoratedCar)
}
