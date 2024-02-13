# 适配器模式
> 适配器适合用于解决新旧系统（或新旧接口）之间的兼容问题，而不建议在一开始就直接使用
- 适配器模式将一个类的接口，转换成客户期望的另一个接口。适配器让原本接口不兼容的类可以合作无间

关键代码:
适配器中持有旧接口对象，并实现新接口

示例：
1. 定义阿里支付接口 type AliPayInterface interface，包含Pay方法
2. 定义微信支付接口 type WeChatPayInterface interface，包含Pay方法
3. 分别定义阿里支付和微信支付的实现类 AliPay 和 WeChatPay struct
4. 定义目标接口 type TargetPayInterface interface，包含DealPay方法
5. 定义TargetPayInterface接口实现类 PayAdapter struct，实现DealPay方法
   - 内部持有AliPayInterface和WeChatPayInterface对象，根据类型分别调用其Pay方法