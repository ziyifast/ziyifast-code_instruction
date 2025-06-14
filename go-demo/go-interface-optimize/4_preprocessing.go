package main

import (
	"fmt"
	"time"
)

// ========== 预处理逻辑 ==========
var precomputed = make(map[int]int)

func precompute() {
	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Millisecond) // 模拟预计算开销
		precomputed[i] = i * i           // 提前算好常用值
	}
}

func usePrecomputed(n int) int {
	return precomputed[n]
}

func runWithPrecompute() {
	startTime := time.Now()

	for i := 0; i < 1000; i++ {
		usePrecomputed(i % 100)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("✅ 启用预处理模式执行完成，总耗时：%v\n", elapsed)
}

// ========== 实时计算逻辑 ==========
func realTimeCompute(n int) int {
	time.Sleep(1 * time.Millisecond) // 模拟实时计算开销
	return n * n
}

func runWithoutPrecompute() {
	startTime := time.Now()

	for i := 0; i < 1000; i++ {
		realTimeCompute(i % 100)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("❌ 无预处理模式执行完成，总耗时：%v\n", elapsed)
}

// ================================
var withPrecomputeElapsed time.Duration
var withoutPrecomputeElapsed time.Duration

func init() {
	// 预先运行一次获取耗时数据
	start := time.Now()
	precompute()
	withPrecomputeElapsed = time.Since(start)

	startNoPre := time.Now()
	for i := 0; i < 1000; i++ {
		realTimeCompute(i % 100)
	}
	withoutPrecomputeElapsed = time.Since(startNoPre)
}

func main() {
	// 执行预处理版本
	runWithPrecompute()

	// 执行非预处理版本
	runWithoutPrecompute()

	// 输出性能对比
	fmt.Println("\n📊 性能对比：")
	fmt.Printf("启用预处理耗时: %v\n", withPrecomputeElapsed)
	fmt.Printf("无预处理耗时: %v\n", withoutPrecomputeElapsed)
	fmt.Printf("优化提升: %.2f 倍\n", float64(withoutPrecomputeElapsed)/float64(withPrecomputeElapsed))
}
