# 访问者模式

> 访问者模式是一种操作一组对象的操作，它的目的是不改变对象的定义，但可以新增不同的访问者来定义新的操作。
> 访问者的核心思想是为了访问比较复杂的数据结构，不去改变原数据结构，而是把对数据的操作抽象出来，在访问的过程中以回调形式在访问者中处理逻辑操作。
> 如果要新增一组操作，那么只需要增加一个新的访问者。


示例：
> 根据不同环境打印不同内容，在生产环境和开发环境中打印处不同的内容

1. 定义IVisitor interface
    - Visit()
2. 定义ProductionVisitor struct，实现IVisitor接口
    - env string
    - func Visit():判断env如果是生产环境，则打印生产环境的内容
3. 定义DevelopmentVisitor struct，实现IVisitor接口
    - env string
    - func Visit():判断env如果是开发环境，则打印开发环境的内容
4. 定义IElement interface接口
    - func Accept(IVisitor)
5. 定义Element struct，实现IElement接口
    - visitors []IVisitor
6. 定义ExampleLog struct（要打印的日志）
   - Element
   - 实现print方法，遍历visitors，调用每个visitor的Visit方法