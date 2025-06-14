package main

import (
	"fmt"
	"sync"
	"time"
)

// ========== 对象池逻辑 ==========
var objectPool = sync.Pool{
	New: func() interface{} {
		time.Sleep(10 * time.Millisecond) // 模拟首次创建对象的耗时
		return &SomeObject{
			Data: make([]byte, 1024),
		}
	},
}

type SomeObject struct {
	Data []byte
}

// 预分配一定数量的对象到池中
func preAllocateToPool(count int) {
	for i := 0; i < count; i++ {
		objectPool.Put(objectPool.New())
	}
}

func usePool() {
	obj := objectPool.Get().(*SomeObject)
	time.Sleep(100 * time.Microsecond) // 模拟使用对象
	objectPool.Put(obj)
}

func runWithPool(iterations int) time.Duration {
	const preAllocCount = 1000
	preAllocateToPool(preAllocCount)

	startTime := time.Now()

	for i := 0; i < iterations; i++ {
		usePool()
	}

	return time.Since(startTime)
}

// ========== 非池化逻辑 ==========
func createNewObject() *SomeObject {
	time.Sleep(5 * time.Millisecond) // 模拟创建client耗时
	return &SomeObject{
		Data: make([]byte, 1024),
	}
}

func runWithoutPool(iterations int) time.Duration {
	startTime := time.Now()

	for i := 0; i < iterations; i++ {
		obj := createNewObject()
		time.Sleep(100 * time.Microsecond) // 模拟使用对象
		_ = obj                            // 不再使用，等待GC回收
	}

	return time.Since(startTime)
}

// ================================
func main() {
	const iterations = 1000

	// 执行对象池版本
	withPoolElapsed := runWithPool(iterations)
	fmt.Printf("✅ 启用对象池模式执行完成，总耗时：%v\n", withPoolElapsed)

	// 执行非对象池版本
	withoutPoolElapsed := runWithoutPool(iterations)
	fmt.Printf("❌ 无对象池模式执行完成，总耗时：%v\n", withoutPoolElapsed)

	// 输出性能对比
	fmt.Println("\n📊 性能对比：")
	fmt.Printf("启用对象池耗时: %v\n", withPoolElapsed)
	fmt.Printf("无对象池耗时: %v\n", withoutPoolElapsed)
	fmt.Printf("优化提升: %.2f 倍\n", float64(withoutPoolElapsed)/float64(withPoolElapsed))
}
