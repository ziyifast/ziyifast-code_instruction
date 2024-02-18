# 职责链模式

> 职责链模式是一种行为设计模式，定义了一系列对象，每个对象可以选择处理某个请求，也可以将该请求传给链中的下一个对象。


模式结构（职责链模式包含以下角色）：
- 抽象处理器（Handler）：定义出一个处理请求的接口。 
- 具体处理器（ConcreteHandler）：实现抽象处理器的接口，处理它所负责的请求。如果不处理该请求，则把请求转发给它的后继者。 
- 客户端（Client）：创建处理器对象，并将请求发送到某个处理器。

示例：
> 演示国家刑事案件处理流程，下级先处理，如果处理不了，交给下一个

1. 定义Handler interface接口
   - SetNext(handler Handler) //设置下一个处理器
   - Handle(request int) //处理请求
2. 定义type TownHandler struct 
    - NextHandler Handler //它的下一个处理器
    - Handle(request int) ：如果案件级别高于20，就交给下一handler处理
    - func SetNext(handler Handler) 
3. 定义type CityHandler struct
    - NextHandler Handler //它的下一个处理器
    - Handle(request int) ：如果案件级别高于100，就交给下一handler处理
    - func SetNext(handler Handler)
4. 定义type ProvinceHandler struct
    - NextHandler Handler //它的下一个处理器
    - Handle(request int) 
    - func SetNext(handler Handler)