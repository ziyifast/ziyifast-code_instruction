# 解释器模式

> 解释器模式（Interpreter Pattern）是一种行为设计模式，它定义了一种语言，用于解释一些特定的领域问题。
> 在该模式中，将语言中的元素映射到类中，并定义它们之间的关系。然后，可以使用这些类来解释表达式，以解决特定的问题。

应用场景：
1. 处理配置文件：json、yaml
2. 模板引擎：模板引擎处理模板和一组变量以产生输出。模板是DSL的一个例子，可以使用Interpreter来解析和处理模板。
3. 数学表达式计算器

解释器模式中的关键组件：

表达式接口：表示抽象语法树的元素并定义解释表达式的方法。
具体表达式：实现表达式接口的结构，表示语言语法的各种规则或元素。
上下文对象：用于保存解释过程中所需的任何必要信息或状态。
Parser 或 Builder：负责根据输入表达式构建抽象语法树的组件。

示例：
> 构建一个计算器
1. 定义Expression interface接口，包含Interpret() int方法
2. 定义NumberExpression struct，实现Interpret() int方法
    - val int
3. 定义AdditionExpression struct，实现Interpret() int方法，返回两数之和
    - left Expression
    - right Expression
    - func Interpret() int {
    }
4. 定义SubtractionExpression struct，实现Interpret() int方法，返回两数之差
    - left Expression
    - right Expression
    - func Interpret() int {
      }
5. 定义Parser struct，存储输入表达式exp、index解析索引位置、prev 上一个表达式
    - exp   []string
    - index int
    - prev  Expression
    - func Parse() Expression {
      }
    - func newAdditionExpression() Expression{
      }
    - func newSubtractionExpression() Expression{
      }