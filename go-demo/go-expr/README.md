# 【go类库分享】Go expr 通用表达式引擎
> 官方教程：https://expr-lang.org/docs/language-definition<br/>
官方Github：https://github.com/expr-lang/expr

# 一、介绍
Expr表达式引擎是一个针对Go语言设计的动态配置解决方案，它以简单的语法和强大的性能特性著称。Expr表达式引擎的核心是安全、快速和直观，很适合用于处理诸如访问控制、数据过滤和资源管理等场景。在Go语言中应用Expr，可以极大地提升应用程序处理动态规则的能力。不同于其他语言的解释器或脚本引擎，Expr采用了静态类型检查，并且生成字节码来执行，因此它能同时保证性能和安全性。

# 二、安装

```go
//通过go get直接安装即可
go get github.com/expr-lang/expr
```

# 三、使用
## 基础使用
### ①运行基本表达式
> 在下面例子中，表达式2 + 2被编译成能运行的字节码，然后执行这段字节码并输出结果。
同时下面的例子不包含变量，因此也不用传入环境。

```go
package main

import (
        "fmt"
        "github.com/expr-lang/expr"
)

func main() {
        // 编译一个基础的加法表达式
        program, err := expr.Compile(`2 + 2`)
        if err != nil {
                panic(err)
        }

        // 运行编译后的表达式，并没有传入环境，因为这里不需要使用任何变量
        output, err := expr.Run(program, nil)
        if err != nil {
                panic(err)
        }

        // 打印结果
        fmt.Println(output)  // 输出 4
}
```

### ②运行变量表达式
> 下面我们创建一个包含变量的环境，编写使用这些变量的表达式，编译并运行这个表达式。
在下面例子中，环境env包含了变量apple和banana。表达式apple + banana在编译时会从环境中推断apple和banana的类型，并在运行时使用这些变量的值来评估表达式结果。

```go
package main

import (
    "fmt"
    "github.com/expr-lang/expr"
)

func main() {
    // 创建一个包含变量的环境
    env := map[string]interface{}{
       "apple":  5,
       "banana": 10,
    }

    // 编译一个使用环境中变量的表达式
    program, err := expr.Compile(`apple + banana`, expr.Env(env))
    if err != nil {
       panic(err)
    }

    // 运行表达式
    output, err := expr.Run(program, env)
    if err != nil {
       panic(err)
    }

    // 打印结果
    fmt.Println(output) // 输出 15
}
```


## 语法介绍
> 下面主要是介绍 Expr 表达式引擎内置函数的一部分。通过这些功能强大的函数，可以更加灵活和高效地处理数据和逻辑。更详细的函数列表和使用说明查阅官方函数文档。

> 官方函数文档：https://expr-lang.org/docs/language-definition


### ①字面量和变量
> Expr表达式引擎能够处理常见的数据类型字面量，包括数字、字符串和布尔值。字面量是直接在代码中写出的数据值，比如42、"hello"和true都是字面量


(1)字面量：数字、字符串、布尔值

```go
// (1) 数字
42      // 表示整数 42
3.14    // 表示浮点数 3.14


// (2) 字符串
"hello, world" // 双引号包裹的字符串，支持转义字符
`hello, world` // 反引号包裹的字符串，保持字符串格式不变，不支持转义

// (3)布尔值
true   // 布尔真值
false  // 布尔假值
```


(2)变量：Expr允许在环境中定义变量，然后在表达式中引用这些变量

```go
// (1)表达式定义变量
env := map[string]interface{}{
    "age": 25,
    "name": "Alice",
}

// (2)表达式中引用变量
age > 18  // 检查age是否大于18
name == "Alice"  // 判断name是否等于"Alice"
```

### ②运算符
Expr表达式引擎支持多种运算符，包含数学运算符、逻辑运算符、比较运算符及集合运算符等。

1. 数学和逻辑运算符
> 数学运算符包括加(+)、减(-)、乘(*)、除(/)和取模(%)。逻辑运算符包括逻辑与(&&)、逻辑或(||)和逻辑非(!)

```go
2 + 2 // 计算结果为4
7 % 3 // 结果为1
!true // 结果为false
age >= 18 && name == "Alice" // 检查age是否不小于18且name是否等于"Alice"
```

2. 比较运算符
> 比较运算符有相等(==)、不等(!=)、小于(<)、小于等于(<=)、大于(>)和大于等于(>=),用于比较两个值

```go
age == 25 // 检查age是否等于25
age != 18 // 检查age是否不等于18
age > 20  // 检查age是否大于20
```
3. 集合运算符
> Expr还提供了一些用于操作集合的运算符，如in用于检查元素是否在集合中，集合可以是数组、切片或字典

```go
"user" in ["user", "admin"]  // true，因为"user"在数组中
3 in {1: true, 2: false}     // false，因为3不是字典的键
```
还有一些高级的集合操作函数，比如all、any、one和none，这些函数需要结合匿名函数(lambda)使用：

```go
all(tweets, {.Len <= 240})  // 检查所有tweets的Len字段是否都不超过240
any(tweets, {.Len > 200})   // 检查是否存在tweets的Len字段超过200
```
4. 成员操作符
> 在Expr表达式语言中，成员操作符允许我们访问Go语言中struct的属性。这个特性让Expr可以直接操作复杂数据结构，非常地灵活实用。

```go
// (1) 定义结构体
type User struct {
    Name string
    Age  int
}

// (2)访问结构体变量
env := map[string]interface{}{
    "user": User{Name: "Alice", Age: 25},
}

code := `user.Name`

program, err := expr.Compile(code, expr.Env(env))
if err != nil {
    panic(err)
}

output, err := expr.Run(program, env)
if err != nil {
    panic(err)
}

fmt.Println(output) // 输出: Alice
```
在操作结构体变量时，我们通常会需要判断对应字段值是否为空，这时就需要处理nil的情况：
> 在访问属性时，可能会遇到对象是nil的情况。Expr提供了安全的属性访问，即使在结构体或者嵌套属性为nil的情况下，也不会抛出运行时panic错误。

方法一：使用`?.`操作符引用属性，如果对象为nil则返回nil，而不会报错。
```go
author.User?.Name

// 等价于下面的表达式
author.User != nil ? author.User.Name : nil
```
方法二：`??`操作符，主要用于nil时，返回默认值

```go
author.User?.Name ?? "Anonymous"

// 等价于下面表达式
author.User != nil ? author.User.Name : "Anonymous"
```
### ③函数
Expr支持内置函数和自定义函数，使得表达式更加强大和灵活。
1. 内置函数：内置函数像len、all、none、any等可以直接在表达式中使用
- all：函数 all 可以用来检验集合中的元素是否全部满足给定的条件。它接受两个参数，第一个参数是集合，第二个参数是条件表达式。

```go
// 检查所有 tweets 的 Content 长度是否小于 240
code := `all(tweets, len(.Content) < 240)`
program, err := expr.Compile(code, expr.Env(env))
if err != nil {
    panic(err)
}
```
- any：与 all 类似，any 函数用来检测集合中是否有任一元素满足条件。

```go
// 检查是否有任一 tweet 的 Content 长度大于 240
code := `any(tweets, len(.Content) > 240)`
```
- none：用于检查集合中没有任何元素满足条件。

```go
// 确保没有 tweets 是重复的
code := `none(tweets, .IsRepeated)`
```

```go
// 内置函数示例
program, err := expr.Compile(`all(users, {.Age >= 18})`, expr.Env(env))
if err != nil {
    panic(err)
}

// 注意：这里env需要包含users变量，每个用户都需要有Age属性
output, err := expr.Run(program, env)
fmt.Print(output) // 如果env中所有用户年龄都大于等于18，返回true
```
2. 自定义函数：通过在环境映射env中传递函数定义来创建自定义函数
> 在Expr中使用函数时，我们可以让代码模块化并在表达式中加入复杂逻辑。通过结合变量、运算符和函数。但需要注意，在构建Expr环境并运行表达式时，始终要确保类型安全。

```go
// 自定义函数示例
env := map[string]interface{}{
    "greet": func(name string) string {
        return fmt.Sprintf("Hello, %s!", name)
    },
}

program, err := expr.Compile(`greet("World")`, expr.Env(env))
if err != nil {
    panic(err)
}

output, err := expr.Run(program, env)
fmt.Print(output) // 返回 Hello, World!
```


## 实际生产案例
>比如我们现在有一个需求：电商平台需要根据用户属性（会员等级、地域）和订单信息（金额、商品类目），动态配置促销活动的参与条件和折扣规则，无需修改代码即可更新规则。

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/expr-lang/expr"
    "github.com/expr-lang/expr/vm"
)

// 用户信息
type User struct {
    ID       int
    Name     string
    Level    int    // 会员等级（1-普通, 2-黄金, 3-钻石）
    Region   string // 用户所在地区
    JoinTime time.Time
}

// 订单信息
type Order struct {
    OrderID     string
    Amount      float64 // 订单金额
    Category    string  // 商品类目（electronics, clothing, food）
    CreatedTime time.Time
}

// 促销规则配置
type PromotionRule struct {
    Condition string  // Expr表达式，判断是否满足条件
    Discount  float64 // 折扣比例（0.9表示9折）
}

// 初始化规则引擎环境
func createEnv(user User, order Order) map[string]interface{} {
    return map[string]interface{}{
       "User":  user,
       "Order": order,
       "Now":   time.Now(), // 内置当前时间函数
       // 可添加其他辅助函数，如字符串处理、数学计算等
    }
}

// 编译促销规则条件
func compileRule(rule string) (*vm.Program, error) {
    return expr.Compile(rule, expr.Env(createEnv(User{}, Order{})))
}

// 应用促销规则
func ApplyPromotion(user User, order Order, rule PromotionRule) (bool, float64, error) {
    // 1. 编译规则（生产环境需缓存编译结果）
    program, err := compileRule(rule.Condition)
    if err != nil {
       return false, 0, fmt.Errorf("规则编译失败: %v", err)
    }

    // 2. 创建执行环境
    env := createEnv(user, order)

    // 3. 执行规则判断
    output, err := expr.Run(program, env)
    if err != nil {
       return false, 0, fmt.Errorf("规则执行失败: %v", err)
    }

    // 4. 类型断言判断结果
    conditionMet, ok := output.(bool)
    if !ok {
       return false, 0, fmt.Errorf("规则必须返回布尔值")
    }

    // 5. 返回是否满足条件及折扣
    return conditionMet, rule.Discount, nil
}

func main() {
    // 模拟用户和订单数据
    user := User{
       ID:       1001,
       Name:     "Alice",
       Level:    3,
       Region:   "CN",
       JoinTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
    }
    order := Order{
       OrderID:     "20231020001",
       Amount:      1500.00,
       Category:    "electronics",
       CreatedTime: time.Now(),
    }

    // 从数据库/配置中心读取促销规则（示例）
    rules := []PromotionRule{
       {
          // 规则1：钻石会员且订单金额>1000，享85折
          Condition: `User.Level >= 3 && Order.Amount > 1000 && Order.Category == "electronics"`,
          Discount:  0.85,
       },
       {
          // 规则2：注册超过2年的用户，任意订单享9折
          Condition: `Now.Sub(User.JoinTime).Hours() > 24*365*2`,
          Discount:  0.9,
       },
    }

    // 遍历所有规则，应用最优折扣
    bestDiscount := 1.0 // 默认无折扣
    for _, rule := range rules {
       valid, discount, err := ApplyPromotion(user, order, rule)
       if err != nil {
          log.Printf("规则应用错误: %v", err)
          continue
       }
       if valid && discount < bestDiscount {
          bestDiscount = discount
       }
    }

    // 计算最终价格
    finalPrice := order.Amount * bestDiscount
    fmt.Printf("原价: ¥%.2f\n", order.Amount)
    fmt.Printf("适用折扣: %.0f%%\n", (1-bestDiscount)*100)
    fmt.Printf("最终价格: ¥%.2f\n", finalPrice)
}
```
## 适用场景

|场景特征  | 推荐方案 | 理由 |
|--|--| -- |
| 规则每天调整多次 |表达式引擎  |  避免频繁发版，提升业务敏捷性  |
| 规则复杂且嵌套业务对象 |直接代码  |   复杂逻辑更易维护，编译器辅助类型检查 |
|  需非技术人员配置规则(产品/运营)|表达式引擎  |  降低技术门槛，释放开发资源  |
| 性能敏感（如：>10万QPS） |  直接代码|  避免表达式解析开销影响吞吐量  |
| 多租户定制规则 | 表达式引擎 |  各租户独立配置，互不影响  |


> 还是以上面的电商场景为例，让大家感受expr的好处以及使用场景：
场景：电商促销规则判断
需求：根据用户等级、订单金额、商品类目动态调整折扣。

方案一：表达式引擎（expr）

```go
// 规则配置（存储于数据库）
rules := []PromotionRule{
    {
        Condition: `User.Level >= 3 && Order.Amount > 1000 && Order.Category == "electronics"`,
        Discount:  0.85,
    },
}
// 动态执行
valid, _ := ApplyPromotion(user, order, rule)
```
优势：
- 运营人员可通过管理后台随时新增/修改规则，无需等待版本发布。
- 支持A/B测试：为不同用户组配置不同规则。

劣势：
- 需额外开发规则管理界面和测试工具。


方案二：直接代码判断

```go
func IsPromotionValid(user User, order Order) bool {
    return user.Level >= 3 && 
           order.Amount > 1000 && 
           order.Category == "electronics"
}
```
优势：
- 性能极高，适合每秒数十万次调用的场景。
- 逻辑变更通过代码评审，降低错误风险。

劣势：
- 修改折扣条件需发版，无法快速响应市场活动。



