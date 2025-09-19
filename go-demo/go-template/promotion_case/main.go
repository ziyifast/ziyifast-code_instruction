package main

import (
	"fmt"
	"os"
	"text/template"
)

// 模拟服务数据
func getUserService(userID string) (map[string]interface{}, error) {
	// 模拟API调用
	user := map[string]interface{}{
		"id":    userID,
		"name":  "张三",
		"email": "zhangsan@example.com",
		"level": "VIP",
	}
	return user, nil
}

func getOrderService(userID string) ([]map[string]interface{}, error) {
	// 模拟API调用
	orders := []map[string]interface{}{
		{
			"id":     "order001",
			"amount": 299.99,
			"status": "已完成",
		},
		{
			"id":     "order002",
			"amount": 199.50,
			"status": "处理中",
		},
	}
	return orders, nil
}

func getPaymentService(orderID string) (map[string]interface{}, error) {
	// 模拟API调用
	payment := map[string]interface{}{
		"method":     "信用卡",
		"status":     "已支付",
		"created_at": "2023-01-01 10:30:00",
	}
	return payment, nil
}

func main() {
	//进阶案例：工作流处理
	funcMap := template.FuncMap{
		"getUser":    getUserService,
		"getOrders":  getOrderService,
		"getPayment": getPaymentService,
	}

	tmpl := `
{{$user := getUser "12345"}}
用户信息:
  姓名: {{$user.name}}
  邮箱: {{$user.email}}
  等级: {{$user.level}}

{{$orders := getOrders "12345"}}
订单列表:
{{range $index, $order := $orders}}
  {{add $index 1}}. 订单ID: {{$order.id}}, 金额: {{$order.amount}}, 状态: {{$order.status}}
  
  {{$payment := getPayment $order.id}}
  支付信息:
    支付方式: {{$payment.method}}
    支付状态: {{$payment.status}}
    支付时间: {{$payment.created_at}}
{{else}}
  没有找到订单
{{end}}`

	// 添加辅助函数
	funcMap["add"] = func(a, b int) int { return a + b }

	t := template.Must(template.New("orchestration").Funcs(funcMap).Parse(tmpl))

	err := t.Execute(os.Stdout, nil)
	if err != nil {
		fmt.Printf("执行模板错误: %v\n", err)
	}
}
