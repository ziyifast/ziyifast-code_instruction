# 享元模式
> 享元模式从对象中剥离出不发生改变且多个实例需要的重复数据，独立出一个享元，使多个对象共享，从而节省内存以及减少对象数量。

示例：
1. 定义结构体ImageFlyweight struct,用于展示图像信息
2. 定义ImageFlyweightFactory struct工厂，包含：maps map[string]*ImageFlyweight
3. 提供方法：GetImageFlyweight(name string) *ImageFlyweight
    - 先从map中获取，如果map中存在，则直接返回；如果不存在则New一个ImageFlyweight，存入map中，然后返回
4. 提供方法GetImageFlyweightFactory，用于获取工厂