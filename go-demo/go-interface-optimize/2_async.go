package main

import (
	"fmt"
	"sync"
	"time"
)

// ========== 异步逻辑 ==========
var wg sync.WaitGroup

func asyncTask() {
	defer wg.Done()
	time.Sleep(2 * time.Second) // 模拟耗时任务
}

func runAsyncTasks() {
	startTime := time.Now()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go asyncTask()
	}
	wg.Wait() // 等待所有 goroutine 完成

	elapsed := time.Since(startTime)
	fmt.Printf("✅ 异步模式执行完成，总耗时：%v\n", elapsed)
}

// ========== 同步逻辑 ==========
func syncTask() {
	time.Sleep(2 * time.Second) // 模拟耗时任务,比如调用给用户发消息接口
}

func runSyncTasks() {
	startTime := time.Now()

	for i := 0; i < 5; i++ {
		syncTask()
	}

	elapsed := time.Since(startTime)
	fmt.Printf("❌ 同步模式执行完成，总耗时：%v\n", elapsed)
}

// ================================
func main() {
	// 测试异步执行
	runAsyncTasks()

	// 测试同步执行
	runSyncTasks()

	// 对比结果
	fmt.Println("\n📊 性能对比：")
	fmt.Printf("异步耗时: %v\n", asyncElapsedTime)
	fmt.Printf("同步耗时: %v\n", syncElapsedTime)
	fmt.Printf("优化提升: %.2f 倍\n", float64(syncElapsedTime)/float64(asyncElapsedTime))
}

var asyncElapsedTime time.Duration
var syncElapsedTime time.Duration

func init() {
	// 预先运行一次以获取耗时数据
	startAsync := time.Now()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go asyncTask()
	}
	wg.Wait()
	asyncElapsedTime = time.Since(startAsync)

	startSync := time.Now()
	for i := 0; i < 5; i++ {
		syncTask()
	}
	syncElapsedTime = time.Since(startSync)
}
