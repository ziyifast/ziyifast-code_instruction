package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const maxConcurrency = 3 // 最大并发数

// ========== 并行逻辑 ==========
func doWork(id int, wg *sync.WaitGroup, ctx context.Context, tokenChan chan struct{}) {
	defer wg.Done()

	select {
	case tokenChan <- struct{}{}: // 获取令牌
	case <-ctx.Done(): // 上下文取消
		fmt.Printf("💤 任务 %d 被取消\n", id)
		return
	}

	defer func() {
		<-tokenChan // 释放令牌
	}()

	// 模拟耗时操作
	time.Sleep(100 * time.Millisecond)
}

func runParallel(ctx context.Context, totalTasks int) time.Duration {
	var wg sync.WaitGroup
	startTime := time.Now()

	tokenChan := make(chan struct{}, maxConcurrency) // 控制并发的 channel
	wg.Add(totalTasks)

	for i := 0; i < totalTasks; i++ {
		go func(i int) {
			doWork(i, &wg, ctx, tokenChan)
		}(i)
	}

	wg.Wait()
	return time.Since(startTime)
}

// ========== 串行逻辑 ==========
func runSerial(totalTasks int) time.Duration {
	startTime := time.Now()

	for i := 0; i < totalTasks; i++ {
		time.Sleep(100 * time.Millisecond)
	}

	return time.Since(startTime)
}

// ================================
func main() {
	const totalTasks = 10
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 执行并行版本（带并发限制）
	parallelElapsed := runParallel(ctx, totalTasks)
	fmt.Printf("✅ 并行模式执行 %d 个任务，总耗时：%v\n", totalTasks, parallelElapsed)

	// 执行串行版本
	serialElapsed := runSerial(totalTasks)
	fmt.Printf("❌ 串行模式执行 %d 个任务，总耗时：%v\n", totalTasks, serialElapsed)

	// 输出性能对比
	fmt.Println("\n📊 性能对比：")
	fmt.Printf("并行耗时: %v\n", parallelElapsed)
	fmt.Printf("串行耗时: %v\n", serialElapsed)
	fmt.Printf("优化提升: %.2f 倍\n", float64(serialElapsed)/float64(parallelElapsed))
}
