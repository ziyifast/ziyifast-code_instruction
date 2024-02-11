# 工厂方法模式 factory method

工厂方法模式使用子类的方式延迟生成对象到子类中实现。

Go 中不存在继承 所以使用匿名组合来实现

示例步骤：
1. 定义接口type operator interface
    - 参数a
    - 参数b
    - result：具体业务方法
2. 定义type BaseFactory struct：提供方法，用于设置a、b参数
    - 参数a
    - 参数b
3. 根据不同操作，定义不同工厂类（addFactory、minusFactory）
   - addFactory实现operator的result：a+b
   - minusFactory实现operator的result：a-b
4. addFactory、minusFactory分别提供Create方法

简单工厂：唯一工厂类，一个产品抽象类，工厂类的创建方法依据入参判断并创建具体产品对象。
工厂方法：多个工厂类，一个产品抽象类，利用多态创建不同的产品对象，避免了大量的if-else判断。
抽象工厂：多个工厂类，多个产品抽象类，产品子类分组，同一个工厂实现类创建同组中的不同产品，减少了工厂子类的数量。