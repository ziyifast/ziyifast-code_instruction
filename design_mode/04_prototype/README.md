# 原型模式

原型模式使对象能复制自身，并且暴露到接口中，使客户端面向接口编程时，不知道接口实际对象的情况下生成新的对象。

原型模式配合原型管理器使用，使得客户端在不知道具体类的情况下，通过接口管理器得到新的实例，并且包含部分预设定配置。

示例步骤：
1. 定义接口type Prototype interface
   - 包含clone方法：Clone() Prototype
2. 定义结构体type ConcretePrototype struct
    - 实现clone方法：Clone() Prototype