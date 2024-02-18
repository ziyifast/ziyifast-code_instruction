# 观察者模式

观察者模式用于触发联动。

一个对象的改变会触发其它观察者的相关动作，而此对象无需关心连动对象的具体实现。
关键：被观察者持有了集合存放观察者 (收通知的为观察者)；类比消息队列的发布订阅，你订阅了此类消息，当有消息来时，我就通知你

示例：
1. 定义Customer interface，包含update() 方法
2. 定义CustomerA、CustomerB结构体，实现update()方法
3. 定义newsOffice struct，包含customers字段，用于存储观察者，后续有事件发生时通知观察者
    - customers []Customer
    - func addCustomer(customer Customer)
    - func removeCustomer(customer Customer)
    - func newsComing() 新报纸到来，调用notifyAll方法通知所有观察者
    - func notifyAll()
