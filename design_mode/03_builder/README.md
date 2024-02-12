# 建造者模式
将build一个物品拆分为几个部分
示例步骤：
1. 定义接口type goodsBuilder interface
    - setName
    - setPrice
    - setCount
    - *Goods
2. 定义具体实现结构体type ConcreteBuilder struct
   - 实现goodsBuilder接口
3. 提供NewGoodsBuilder接口，返回ConcreteBuilder实现类