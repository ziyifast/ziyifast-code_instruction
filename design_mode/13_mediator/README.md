# 中介者模式
中介者模式封装对象之间互交，使依赖变的简单，并且使复杂互交简单化，封装在中介者中。

例子中的中介者使用单例模式生成中介者。

中介者的change使用switch判断类型。

> 使用场景：
> 1. 当一些对象和其他对象紧密耦合以致难以对其进行修改时，可使用中介者模式。 
> 2. 当组件因过于依赖其他组件而无法在不同应用中复用时，可使用中介者模式。 
> 3. 如果为了能在不同情景下复用一些基本行为，导致你需要被迫创建大量组件子类时，可使用中介者模式。
> 

示例：
1. 定义MessageMediator interface中介者接口，包含sendMessage(msg string, user User)、receiveMessage() string方法
2. 定义具体中介者实现类，ChatRoom struct，属性包含Message string
   - 实现sendMessage方法，将消息发送给用户
3. 定义type User struct
   - name     string
   - mediator Mediator
4. 用户通过mediator发送消息，通过mediator（chatroom）接受消息