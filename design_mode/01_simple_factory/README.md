# 简单工厂模式 simple factory
go 语言没有构造函数，所以我们一般是通过 NewXXX 函数来初始化相关类。 NewXXX 函数返回接口时就是简单工厂模式，也就是说 Golang 的一般推荐做法就是简单工厂。

在这个 simplefactory 包中只有API 接口和 NewAPI 函数为包外可见，封装了实现细节。