# 策略模式
定义一系列算法，让这些算法在运行时可以互换，使得分离算法，符合开闭原则。

示例：
1. 定义PaymentStrategy interface 接口
    - Pay(amount float64) error
2. 定义实现类，CreditCardPaymentStrategy、CashPaymentStrategy
    - 分别实现Pay(amount float64) error方法
3. 定义上下文对象类PaymentContext
    - amount float64
    - strategy PaymentStrategy
4. PaymentContext提供Pay()方法，调用strategy.Pay(amount)方法
   - p.strategy.Pay(p.amount)
5. 提供NewPaymentContext方法，返回*PaymentContext对象
6. 根据不同策略，调用NewPaymentContext方法，返回不同PaymentContext对象，实现不同的支付
