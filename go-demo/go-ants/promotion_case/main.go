package main

import (
	"context"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// OrderProcessor 订单处理器
type OrderProcessor struct {
	// 不同业务类型的专用协程池
	orderValidationPool *ants.Pool // 订单验证池
	paymentProcessPool  *ants.Pool // 支付处理池
	inventoryUpdatePool *ants.Pool // 库存更新池
	emailNotifyPool     *ants.Pool // 邮件通知池

	// 统计数据
	totalProcessed int64
	successCount   int64
	errorCount     int64
}

// Order 订单结构
type Order struct {
	ID          string
	UserID      int64
	Products    []ProductItem
	TotalAmount float64
	Status      string
	CreatedAt   time.Time
}

// ProductItem 商品项
type ProductItem struct {
	ProductID int64
	Quantity  int
	Price     float64
}

// NewOrderProcessor 创建订单处理器
func NewOrderProcessor() *OrderProcessor {
	return &OrderProcessor{
		// 根据业务特点分配不同容量的协程池
		orderValidationPool: createPoolWithConfig(500), // 订单验证较轻量
		paymentProcessPool:  createPoolWithConfig(200), // 支付涉及外部系统，限制并发
		inventoryUpdatePool: createPoolWithConfig(300), // 库存更新需要控制
		emailNotifyPool:     createPoolWithConfig(100), // 邮件通知可以较低优先级
	}
}

// createPoolWithConfig 创建带配置的协程池.具体配置结合自身业务场景
func createPoolWithConfig(size int) *ants.Pool {
	pool, err := ants.NewPool(size,
		ants.WithExpiryDuration(30*time.Second), // 协程空闲30秒后回收
		ants.WithNonblocking(false))             // 阻塞模式，排队等待而非拒绝
	if err != nil {
		log.Fatal("创建协程池失败:", err)
	}
	return pool
}

// ProcessOrder 处理订单
func (op *OrderProcessor) ProcessOrder(order *Order) error {
	atomic.AddInt64(&op.totalProcessed, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var errors []error
	var mu sync.Mutex

	// 1. 订单验证 - 使用专用验证池
	wg.Add(1)
	err := op.orderValidationPool.Submit(func() {
		defer wg.Done()
		if err := op.validateOrder(ctx, order); err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("验证失败: %v", err))
			mu.Unlock()
		}
	})
	if err != nil {
		return fmt.Errorf("验证任务提交失败: %v", err)
	}

	// 2. 支付处理 - 使用专用支付池
	wg.Add(1)
	err = op.paymentProcessPool.Submit(func() {
		defer wg.Done()
		if err := op.processPayment(ctx, order); err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("支付失败: %v", err))
			mu.Unlock()
		}
	})
	if err != nil {
		return fmt.Errorf("支付任务提交失败: %v", err)
	}

	// 等待关键步骤完成
	wg.Wait()

	// 检查是否有错误
	if len(errors) > 0 {
		atomic.AddInt64(&op.errorCount, 1)
		order.Status = "failed"
		return errors[0] // 返回第一个错误
	}

	// 3. 异步更新库存 - 使用库存池
	go func() {
		_ = op.inventoryUpdatePool.Submit(func() {
			op.updateInventory(ctx, order)
		})
	}()

	// 4. 异步发送邮件 - 使用邮件池
	go func() {
		_ = op.emailNotifyPool.Submit(func() {
			op.sendEmailNotification(ctx, order)
		})
	}()

	order.Status = "completed"
	atomic.AddInt64(&op.successCount, 1)
	return nil
}

// validateOrder 订单验证
func (op *OrderProcessor) validateOrder(ctx context.Context, order *Order) error {
	// 模拟数据库查询和业务验证
	time.Sleep(time.Duration(rand.Intn(50)+10) * time.Millisecond)

	// 模拟验证逻辑
	if len(order.Products) == 0 {
		return fmt.Errorf("订单商品不能为空")
	}

	for _, item := range order.Products {
		if item.Quantity <= 0 {
			return fmt.Errorf("商品数量必须大于0")
		}
	}

	log.Printf("订单%s验证通过", order.ID)
	return nil
}

// processPayment 处理支付
func (op *OrderProcessor) processPayment(ctx context.Context, order *Order) error {
	// 模拟调用第三方支付接口
	time.Sleep(time.Duration(rand.Intn(200)+50) * time.Millisecond)

	// 模拟支付成功率
	if rand.Float32() < 0.02 { // 2%失败率
		return fmt.Errorf("支付网关超时")
	}

	log.Printf("订单%s支付成功，金额: %.2f", order.ID, order.TotalAmount)
	return nil
}

// updateInventory 更新库存
func (op *OrderProcessor) updateInventory(ctx context.Context, order *Order) {
	// 模拟库存系统更新
	time.Sleep(time.Duration(rand.Intn(100)+20) * time.Millisecond)
	log.Printf("订单%s库存更新完成", order.ID)
}

// sendEmailNotification 发送邮件通知
func (op *OrderProcessor) sendEmailNotification(ctx context.Context, order *Order) {
	// 模拟邮件发送
	time.Sleep(time.Duration(rand.Intn(300)+50) * time.Millisecond)
	log.Printf("订单%s邮件通知发送完成", order.ID)
}

// GetStats 获取处理统计
func (op *OrderProcessor) GetStats() map[string]int64 {
	return map[string]int64{
		"total_processed": atomic.LoadInt64(&op.totalProcessed),
		"success_count":   atomic.LoadInt64(&op.successCount),
		"error_count":     atomic.LoadInt64(&op.errorCount),
	}
}

// GetPoolStats 获取协程池状态
func (op *OrderProcessor) GetPoolStats() map[string]map[string]int {
	return map[string]map[string]int{
		"validation_pool": {
			"capacity": op.orderValidationPool.Cap(),
			"running":  op.orderValidationPool.Running(),
		},
		"payment_pool": {
			"capacity": op.paymentProcessPool.Cap(),
			"running":  op.paymentProcessPool.Running(),
		},
		"inventory_pool": {
			"capacity": op.inventoryUpdatePool.Cap(),
			"running":  op.inventoryUpdatePool.Running(),
		},
		"email_pool": {
			"capacity": op.emailNotifyPool.Cap(),
			"running":  op.emailNotifyPool.Running(),
		},
	}
}

// Close 关闭所有协程池
func (op *OrderProcessor) Close() {
	op.orderValidationPool.Release()
	op.paymentProcessPool.Release()
	op.inventoryUpdatePool.Release()
	op.emailNotifyPool.Release()
}

// generateTestOrders 生成测试订单
func generateTestOrders(count int) []*Order {
	orders := make([]*Order, count)
	for i := 0; i < count; i++ {
		orders[i] = &Order{
			ID:     fmt.Sprintf("ORDER_%06d", i),
			UserID: int64(rand.Intn(10000) + 1),
			Products: []ProductItem{
				{ProductID: int64(rand.Intn(1000) + 1), Quantity: rand.Intn(5) + 1, Price: float64(rand.Intn(1000))/10 + 10},
			},
			TotalAmount: float64(rand.Intn(5000))/10 + 50,
			CreatedAt:   time.Now(),
		}
	}
	return orders
}

func main() {
	fmt.Println("=== 电商订单处理系统 - ants企业级应用 ===")
	fmt.Printf("CPU核心数: %d\n", runtime.NumCPU())
	fmt.Printf("初始goroutine数: %d\n\n", runtime.NumGoroutine())

	// 创建订单处理器
	processor := NewOrderProcessor()
	defer processor.Close()

	// 生成测试订单
	testOrders := generateTestOrders(5000)
	fmt.Printf("生成测试订单: %d个\n", len(testOrders))

	startTime := time.Now()

	// 模拟高并发订单处理
	var wg sync.WaitGroup

	// 使用多个goroutine并发提交订单处理任务
	workerCount := 100
	ordersPerWorker := len(testOrders) / workerCount

	for worker := 0; worker < workerCount; worker++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			startIdx := workerID * ordersPerWorker
			endIdx := startIdx + ordersPerWorker
			if workerID == workerCount-1 {
				endIdx = len(testOrders) // 处理余数
			}

			for i := startIdx; i < endIdx; i++ {
				err := processor.ProcessOrder(testOrders[i])
				if err != nil {
					log.Printf("Worker%d: 订单%s处理失败: %v", workerID, testOrders[i].ID, err)
					//todo 记录错误，进行补偿/上报/其他处理，此处为了代码简单，不做其他逻辑
				}

				// 模拟订单到达间隔
				if i%100 == 0 {
					time.Sleep(time.Millisecond)
				}
			}
		}(worker)
	}

	wg.Wait()

	processingTime := time.Since(startTime)

	// 输出结果
	fmt.Println("\n=== 处理结果统计 ===")
	stats := processor.GetStats()
	fmt.Printf("总处理时间: %v\n", processingTime)
	fmt.Printf("处理订单总数: %d\n", stats["total_processed"])
	fmt.Printf("成功处理: %d\n", stats["success_count"])
	fmt.Printf("处理失败: %d\n", stats["error_count"])
	fmt.Printf("成功率: %.2f%%\n",
		float64(stats["success_count"])/float64(stats["total_processed"])*100)
	fmt.Printf("处理吞吐量: %.2f 订单/秒\n",
		float64(stats["total_processed"])/processingTime.Seconds())

	// 协程池状态
	fmt.Println("\n=== 协程池状态 ===")
	poolStats := processor.GetPoolStats()
	for poolName, stat := range poolStats {
		fmt.Printf("%s: 容量=%d, 运行中=%d\n",
			poolName, stat["capacity"], stat["running"])
	}

	// 内存使用情况
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\n内存使用: %d MB\n", m.Alloc/1024/1024)
	fmt.Printf("最终goroutine数: %d\n", runtime.NumGoroutine())
}
