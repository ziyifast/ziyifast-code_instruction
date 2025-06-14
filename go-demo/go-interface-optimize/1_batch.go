package main

import (
	"fmt"
	"time"
)

const batchSize = 1000 // 批量大小

var buffer []int // 模拟待插入数据的缓冲区

// ========== 批处理逻辑 ==========
func batchInsert(data int) {
	buffer = append(buffer, data)
	if len(buffer) >= batchSize {
		flushBatch()
	}
}

func flushBatch() {
	if len(buffer) == 0 {
		return
	}
	time.Sleep(1 * time.Millisecond) // 模拟IO延迟
	buffer = buffer[:0]              // 清空缓冲区
}

// ========== 非批处理逻辑 ==========
func singleInsert(data int) {
	time.Sleep(1 * time.Millisecond) // 模拟每次插入都发生IO
}

// ================================
func main() {
	const totalItems = 2000

	// 测试批处理耗时
	startBatch := time.Now()
	buffer = nil // 重置buffer
	for i := 0; i < totalItems; i++ {
		batchInsert(i)
	}
	flushBatch() // 处理剩余数据
	elapsedBatch := time.Since(startBatch)
	fmt.Printf("✅ 批处理模式共 %d 条数据，总耗时：%v\n", totalItems, elapsedBatch)

	// 测试非批处理耗时
	startSingle := time.Now()
	for i := 0; i < totalItems; i++ {
		singleInsert(i)
	}
	elapsedSingle := time.Since(startSingle)
	fmt.Printf("❌ 非批处理模式共 %d 条数据，总耗时：%v\n", totalItems, elapsedSingle)

	// 对比结果
	fmt.Println("\n📊 性能对比：")
	fmt.Printf("批处理耗时: %v\n", elapsedBatch)
	fmt.Printf("非批处理耗时: %v\n", elapsedSingle)
	fmt.Printf("优化提升: %.2f 倍\n", float64(elapsedSingle)/float64(elapsedBatch))
}
