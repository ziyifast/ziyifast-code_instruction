package main

import (
	"fmt"
	"log"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

/*
	实际场景：
		电商平台需要根据用户属性（会员等级、地域）和订单信息（金额、商品类目），
		动态配置促销活动的参与条件和折扣规则，无需修改代码即可更新规则
*/

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
