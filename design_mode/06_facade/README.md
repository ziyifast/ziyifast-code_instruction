# 外观模式
> 外观模式也叫门面模式，是一种结构型设计模式，它提供了一个统一的接口来访问子系统中的一组接口。这种模式通过定义一个高层接口来隐藏子系统的复杂性，使子系统更容易使用。 
> 在Go语言中，我们可以使用结构体和接口来实现外观模式。

示例：
1. 定义audioFixer，实现fixer
2. 定义videoFixer，实现fixer
3. 定义AudioAndVideoFixer，组合audioFixer和videoFixer