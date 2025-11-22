package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

var wg sync.WaitGroup
var sum int64
var mu sync.Mutex

func main() {

	//API_Demo01()
	//BasicBatch_Demo02()
	//PromotionBatch_Demo03()
	//PanicHandler_Demo04()

}

// PanicHandler_Demo04 错误处理
func PanicHandler_Demo04() {
	panicHandler := func(p interface{}) {
		fmt.Printf("Worker exits from panic: %v\n", p)
		// Log the panic or send alert
	}
	pool, _ := ants.NewPool(10, ants.WithPanicHandler(panicHandler))
	defer pool.Release()
	_ = pool.Submit(func() {
		panic("panic error")
	})
	time.Sleep(time.Second)
}

// PromotionBatch_Demo03 批量同类任务处理
func PromotionBatch_Demo03() {
	incr := func(i any) {
		mu.Lock()
		sum += i.(int64)
		mu.Unlock()
		wg.Done()
	}
	p, _ := ants.NewPoolWithFunc(1000, incr)
	defer p.Release()
	// 快速处理大量相似任务
	for i := 0; i < 50000; i++ {
		wg.Add(1)
		_ = p.Invoke(int64(i))
	}
	wg.Wait() // 等待所有任务完成
	fmt.Printf("Final result: %d\n", sum)
}

// BasicBatch_Demo02 批量处理任务
func BasicBatch_Demo02() {
	//处理不同类型(handler中干不同事情)的批量任务
	task1 := func() {
		defer wg.Done()
		fmt.Println("hello...")
	}
	task2 := func() {
		defer wg.Done()
		fmt.Println("world...")
	}
	pool, _ := ants.NewPool(10)
	defer pool.Release()
	for i := 0; i < 20; i++ {
		wg.Add(1)
		_ = pool.Submit(task1)
	}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		_ = pool.Submit(task2)
	}
	wg.Wait()
	fmt.Println("all done.")
}

// API_Demo01 演示基础API
func API_Demo01() {
	//Case 创建协程池：创建一个最大容量为10的协程池
	pool, _ := ants.NewPool(10)
	defer pool.Release() // 释放协程池
	fmt.Println("pool=", pool)

	//Case 提交任务
	task := func() {
		fmt.Println("hello world...")
	}
	_ = pool.Submit(task)
	time.Sleep(time.Second)

	//Case 创建带预定义函数的工作池。然后通过invoke去调用触发
	handler := func(data interface{}) {
		n := data.(int)
		fmt.Printf("Processing number: %d\n", n)
	}
	p, _ := ants.NewPoolWithFunc(500, handler)
	defer p.Release()
	for i := 0; i < 10; i++ {
		_ = p.Invoke(i)
	}
	time.Sleep(time.Second * 3)

	//Case 动态调整协程池大小
	fmt.Printf("before tune, p.size %d\n", p.Cap())
	p.Tune(20)
	fmt.Printf("after tune, p.size %d\n", p.Cap())
}
