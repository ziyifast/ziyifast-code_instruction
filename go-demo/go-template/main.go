package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

// 服务响应结构
type ServiceResponse struct {
	Data interface{} `json:"data"`
}

// 用户信息
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Level    string `json:"level"`
	Region   string `json:"region"`
	Language string `json:"language"`
}

// 订单信息
type Order struct {
	ID         string            `json:"id"`
	Amount     float64           `json:"amount"`
	Status     string            `json:"status"`
	Currency   string            `json:"currency"`
	Items      []OrderItem       `json:"items"`
	Metadata   map[string]string `json:"metadata"`
	CreatedAt  string            `json:"created_at"`
	CustomerID string            `json:"customer_id"`
}

// 订单项
type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// 支付信息
type Payment struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"order_id"`
	Method      string  `json:"method"`
	Status      string  `json:"status"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	CreatedAt   string  `json:"created_at"`
	Transaction string  `json:"transaction_id"`
}

// 库存信息
type Inventory struct {
	ProductID string `json:"product_id"`
	Stock     int    `json:"stock"`
	Reserved  int    `json:"reserved"`
}

// 通知配置
type NotificationConfig struct {
	Channels []string `json:"channels"`
	Template string   `json:"template"`
}

// 模拟数据
var (
	users = map[string]User{
		"cust_12345": {ID: "cust_12345", Name: "张三", Email: "zhangsan@example.com", Level: "VIP", Region: "CN", Language: "zh-CN"},
		"cust_67890": {ID: "cust_67890", Name: "李四", Email: "lisi@example.com", Level: "普通", Region: "US", Language: "en-US"},
	}

	orders = map[string]Order{
		"ord_001": {
			ID:         "ord_001",
			Amount:     299.99,
			Status:     "completed",
			Currency:   "CNY",
			CreatedAt:  "2023-01-01T10:30:00Z",
			CustomerID: "cust_12345",
			Items: []OrderItem{
				{ProductID: "prod_001", ProductName: "iPhone 14", Quantity: 1, Price: 299.99},
			},
			Metadata: map[string]string{
				"source": "web",
				"promo":  "NEWYEAR2023",
			},
		},
		"ord_002": {
			ID:         "ord_002",
			Amount:     199.50,
			Status:     "processing",
			Currency:   "CNY",
			CreatedAt:  "2023-01-02T14:15:00Z",
			CustomerID: "cust_12345",
			Items: []OrderItem{
				{ProductID: "prod_002", ProductName: "AirPods", Quantity: 1, Price: 199.50},
			},
			Metadata: map[string]string{
				"source": "mobile",
			},
		},
	}

	payments = map[string]Payment{
		"pay_001": {ID: "pay_001", OrderID: "ord_001", Method: "credit_card", Status: "success", Amount: 299.99, Currency: "CNY", CreatedAt: "2023-01-01T10:30:00Z", Transaction: "txn_abc123"},
		"pay_002": {ID: "pay_002", OrderID: "ord_002", Method: "alipay", Status: "pending", Amount: 199.50, Currency: "CNY", CreatedAt: "2023-01-02T14:15:00Z", Transaction: "txn_def456"},
	}

	inventories = map[string]Inventory{
		"prod_001": {ProductID: "prod_001", Stock: 100, Reserved: 10},
		"prod_002": {ProductID: "prod_002", Stock: 50, Reserved: 5},
	}

	notificationConfigs = map[string]NotificationConfig{
		"order_completed":  {Channels: []string{"email", "sms"}, Template: "order_completion"},
		"order_processing": {Channels: []string{"email"}, Template: "order_processing"},
	}
)

// 模拟用户服务
func userHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "缺少用户ID参数", http.StatusBadRequest)
		return
	}

	user, exists := users[userID]
	if !exists {
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}

	response := ServiceResponse{Data: user}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 模拟订单服务
func orderHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		http.Error(w, "缺少订单ID参数", http.StatusBadRequest)
		return
	}

	order, exists := orders[orderID]
	if !exists {
		http.Error(w, "订单不存在", http.StatusNotFound)
		return
	}

	response := ServiceResponse{Data: order}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 模拟支付服务
func paymentHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		http.Error(w, "缺少订单ID参数", http.StatusBadRequest)
		return
	}

	// 查找与订单关联的支付信息
	var payment *Payment
	for _, p := range payments {
		if p.OrderID == orderID {
			payment = &p
			break
		}
	}

	if payment == nil {
		http.Error(w, "支付信息不存在", http.StatusNotFound)
		return
	}

	response := ServiceResponse{Data: *payment}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 模拟库存服务
func inventoryHandler(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "缺少产品ID参数", http.StatusBadRequest)
		return
	}

	inventory, exists := inventories[productID]
	if !exists {
		http.Error(w, "库存信息不存在", http.StatusNotFound)
		return
	}

	response := ServiceResponse{Data: inventory}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 模拟通知配置服务
func notificationConfigHandler(w http.ResponseWriter, r *http.Request) {
	eventType := r.URL.Query().Get("event_type")
	if eventType == "" {
		http.Error(w, "缺少事件类型参数", http.StatusBadRequest)
		return
	}

	config, exists := notificationConfigs[eventType]
	if !exists {
		// 返回默认配置
		config = NotificationConfig{Channels: []string{"email"}, Template: "default"}
	}

	response := ServiceResponse{Data: config}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 启动模拟服务
func startMockServices() {
	//调用本地http服务，用于模拟调用下游其他业务服务
	http.HandleFunc("/users", userHandler)
	http.HandleFunc("/orders", orderHandler)
	http.HandleFunc("/payments", paymentHandler)
	http.HandleFunc("/inventory", inventoryHandler)
	http.HandleFunc("/notification-config", notificationConfigHandler)

	fmt.Println("模拟服务已启动在 :8080 端口")
	fmt.Println("用户服务: http://localhost:8080/users?id=cust_12345")
	fmt.Println("订单服务: http://localhost:8080/orders?id=ord_001")
	fmt.Println("支付服务: http://localhost:8080/payments?order_id=ord_001")
	fmt.Println("库存服务: http://localhost:8080/inventory?product_id=prod_001")
	fmt.Println("通知配置服务: http://localhost:8080/notification-config?event_type=order_completed")
	fmt.Println()

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// 等待服务启动
	time.Sleep(100 * time.Millisecond)
}

// 自定义函数：调用用户服务
func getUserService(userID string) (map[string]interface{}, error) {
	resp, err := http.Get("http://localhost:8080/users?id=" + userID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 转换为 map[string]interface{}
	userData, ok := result.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("用户数据格式错误")
	}

	return userData, nil
}

// 自定义函数：调用订单服务
func getOrderService(orderID string) (map[string]interface{}, error) {
	resp, err := http.Get("http://localhost:8080/orders?id=" + orderID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 转换为 map[string]interface{}
	orderData, ok := result.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("订单数据格式错误")
	}

	return orderData, nil
}

// 自定义函数：调用支付服务
func getPaymentService(orderID string) (map[string]interface{}, error) {
	resp, err := http.Get("http://localhost:8080/payments?order_id=" + orderID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 转换为 map[string]interface{}
	paymentData, ok := result.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("支付数据格式错误")
	}

	return paymentData, nil
}

// 自定义函数：调用库存服务
func getInventoryService(productID string) (map[string]interface{}, error) {
	resp, err := http.Get("http://localhost:8080/inventory?product_id=" + productID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 转换为 map[string]interface{}
	inventoryData, ok := result.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("库存数据格式错误")
	}

	return inventoryData, nil
}

// 自定义函数：获取通知配置
func getNotificationConfig(eventType string) (map[string]interface{}, error) {
	resp, err := http.Get("http://localhost:8080/notification-config?event_type=" + eventType)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 转换为 map[string]interface{}
	configData, ok := result.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("通知配置数据格式错误")
	}

	return configData, nil
}

// 自定义函数：格式化货币
func formatCurrency(amount float64, currency string) string {
	switch currency {
	case "CNY":
		return fmt.Sprintf("¥%.2f", amount)
	case "USD":
		return fmt.Sprintf("$%.2f", amount)
	default:
		return fmt.Sprintf("%.2f %s", amount, currency)
	}
}

// 自定义函数：获取订单状态描述
func getOrderStatusDescription(status string) string {
	statusMap := map[string]string{
		"completed":  "已完成",
		"processing": "处理中",
		"cancelled":  "已取消",
		"pending":    "待处理",
	}

	if desc, ok := statusMap[status]; ok {
		return desc
	}
	return status
}

// 辅助函数：加法
func add(a, b int) int {
	return a + b
}

// 辅助函数：乘法
func multiply(a, b float64) float64 {
	return a * b
}

// 辅助函数：连接字符串
func join(sep string, s []interface{}) string {
	if len(s) == 0 {
		return ""
	}
	result := fmt.Sprintf("%v", s[0])
	for _, item := range s[1:] {
		result += sep + fmt.Sprintf("%v", item)
	}
	return result
}

// 辅助函数：转换为JSON
func toJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func main() {
	// 启动模拟服务
	startMockServices()

	// 定义所有函数
	funcMap := template.FuncMap{
		"getUser":               getUserService,
		"getOrder":              getOrderService,
		"getPayment":            getPaymentService,
		"getInventory":          getInventoryService,
		"getNotificationConfig": getNotificationConfig,
		"formatCurrency":        formatCurrency,
		"getOrderStatus":        getOrderStatusDescription,
		"add":                   add,
		"multiply":              multiply,
		"join":                  join,
		"toJson":                toJson,
	}

	// 创建复杂的业务流程编排模板
	tmplText := `
{{- /* 获取订单信息 */ -}}
{{$order := getOrder "ord_001"}}
{{$customer := getUser $order.customer_id}}

订单处理报告
================================
订单ID: {{$order.id}}
订单状态: {{getOrderStatus $order.status}}
创建时间: {{$order.created_at}}
订单金额: {{formatCurrency $order.amount $order.currency}}

客户信息:
  姓名: {{$customer.name}}
  邮箱: {{$customer.email}}
  等级: {{$customer.level}}
  地区: {{$customer.region}}

订单项目:
{{range $index, $item := $order.items}}
  {{$index | add 1}}. {{$item.product_name}}
     产品ID: {{$item.product_id}}
     数量: {{$item.quantity}}
     单价: {{formatCurrency $item.price $order.currency}}
     小计: {{formatCurrency (multiply $item.quantity $item.price) $order.currency}}
{{- end -}}

{{/* 获取支付信息 */}}
{{$payment := getPayment $order.id}}
支付信息:
  支付ID: {{$payment.id}}
  支付方式: {{$payment.method}}
  支付状态: {{$payment.status}}
  交易号: {{$payment.transaction_id}}
  支付时间: {{$payment.created_at}}

{{/* 检查库存信息 */}}
{{$inventoryIssues := false}}
库存检查:
{{- range $item := $order.items -}}
  {{$inventory := getInventory $item.product_id}}
  产品 {{$item.product_name}}:
    可用库存: {{$inventory.stock}}
    已预留: {{$inventory.reserved}}
    {{if lt $inventory.stock $item.quantity}}
    警告: 库存不足！
    {{$inventoryIssues = true}}
    {{end}}
{{- end -}}

{{/* 获取通知配置 */}}
{{$notificationConfig := getNotificationConfig (printf "order_%s" $order.status)}}
通知配置:
  渠道: {{join ", " $notificationConfig.channels}}
  模板: {{- $notificationConfig.template -}}

{{/* 业务决策 */}}
{{if eq $order.status "completed"}}
  {{if eq $payment.status "success"}}
    {{if not $inventoryIssues}}
处理结论: 订单处理成功，可以进行发货准备
    {{else}}
处理结论: 订单支付成功但库存不足，请人工处理
    {{end}}
  {{else}}
处理结论: 订单状态异常，请检查支付状态
  {{end}}
{{else if eq $order.status "processing"}}
处理结论: 订单正在处理中，等待后续状态更新
{{else}}
处理结论: 订单状态{{$order.status}}，按相应流程处理
{{- end -}}

{{/* 生成下游服务调用 */}}
下游服务调用:
{{if eq $order.status "completed"}}
  1. 发货服务: prepareShipment({{toJson $order}})
  2. 通知服务: sendNotification({{toJson $notificationConfig}}, {{$customer.email}})
  3. 积分服务: addPoints({{$customer.id}}, {{$order.amount}})
{{end}}
`

	// 创建模板并添加所有函数
	tmpl := template.Must(template.New("business_process").Funcs(funcMap).Parse(tmplText))

	// 执行模板
	fmt.Println("开始执行业务流程编排...")
	fmt.Println("========================")

	err := tmpl.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatal("执行模板错误: ", err)
	}
}
