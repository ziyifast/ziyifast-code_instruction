# 迭代器模式 
迭代器模式用于使用相同方式送代不同类型集合或者隐藏集合类型的具体实现。

可以使用迭代器模式使遍历同时应用送代策略，如请求新对象、过滤、处理对象等。

示例：
1. 定义Iterator interface
   - Next() interface{}
   - HasNext() bool
2. 定义NumberIterator struct结构体，实现Iterator接口
3. 定义Numbers struct
   - numbers []int
   - func Iterator()，用于获取number的迭代器
    