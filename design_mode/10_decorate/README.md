# 装饰模式
装饰模式使用对象组合的方式动态改变或增加对象行为。

Go语言借助于匿名组合和非入侵式接口可以很方便实现装饰模式。

使用匿名组合，在装饰器中不必显式定义转调原对象方法。

示例：
1. 定义PriceDecorator interface接口
   - 包含DecoratePrice(c Car) Car方法,用于增加车的price
2. 定义ExtraPriceDecorator结构体,实现PriceDecorator接口
   - 包含ExtraPrice字段,用于给传入的car增加额外的price，最后返回car