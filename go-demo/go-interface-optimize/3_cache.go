package main

import (
	"fmt"
	"time"
)

// ========== 缓存逻辑 ==========
var cache = make(map[int]int)

func computeWithCache(n int) int {
	if result, found := cache[n]; found {
		return result
	}
	time.Sleep(1 * time.Millisecond) // 模拟耗时计算
	result := n * n
	cache[n] = result
	return result
}

func runWithCache() {
	startTime := time.Now()

	for i := 0; i < 1000; i++ {
		computeWithCache(i % 100) // 频繁访问少量数据
	}

	elapsed := time.Since(startTime)
	fmt.Printf("✅ 启用缓存模式执行完成，总耗时：%v\n", elapsed)
}

// ========== 无缓存逻辑 ==========
func computeNoCache(n int) int {
	time.Sleep(1 * time.Millisecond) // 模拟耗时计算
	return n * n
}

func runWithoutCache() {
	startTime := time.Now()

	for i := 0; i < 1000; i++ {
		computeNoCache(i % 100) // 重复计算相同输入
	}

	elapsed := time.Since(startTime)
	fmt.Printf("❌ 无缓存模式执行完成，总耗时：%v\n", elapsed)
}

// ================================
func main() {
	// 测试启用缓存
	runWithCache()

	// 测试无缓存
	runWithoutCache()

	// 对比结果
	fmt.Println("\n📊 性能对比：")
	fmt.Printf("启用缓存耗时: %v\n", cachedElapsedTime)
	fmt.Printf("无缓存耗时: %v\n", noCacheElapsedTime)
	fmt.Printf("优化提升: %.2f 倍\n", float64(noCacheElapsedTime)/float64(cachedElapsedTime))
}

var cachedElapsedTime time.Duration
var noCacheElapsedTime time.Duration

func init() {
	// 预先运行一次以获取耗时数据
	startCached := time.Now()
	for i := 0; i < 1000; i++ {
		computeWithCache(i % 100)
	}
	cachedElapsedTime = time.Since(startCached)

	startNoCache := time.Now()
	for i := 0; i < 1000; i++ {
		computeNoCache(i % 100)
	}
	noCacheElapsedTime = time.Since(startNoCache)
}
