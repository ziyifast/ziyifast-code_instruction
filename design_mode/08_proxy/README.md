# 代理模式
> 代理模式用于延迟处理操作或者在进行实际操作前后进行其它处理。


示例：
1. 定义PaymentService interface接口，包含pay方法
2. 定义Alipay struct，实现PaymentService interface接口
3. 定义paymentProxy struct，包含realPay（实现了PaymentService interface的结构体）
4. paymentProxy实现pay方法，调用realPay的pay方法之前和之后分别来添加一些额外的操作