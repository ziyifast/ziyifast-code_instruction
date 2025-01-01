package main

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	//http.HandleFunc("/", handler)
	//fmt.Println("Server started at :8080")
	//http.ListenAndServe(":8080", nil)

	//Demo_CurrentHandleWorks()
	//Demo_CurrentHandleWithHTTP()
	Demo_HandleWorkWithWait()
}

// Demo_HandleWorkWithWait 模拟等待耗时任务处理
func Demo_HandleWorkWithWait() {
	// 创建一个限速器，每3秒允许1个事件
	limiter := rate.NewLimiter(rate.Every(3*time.Second), 1)

	// 模拟10次对资源的访问
	for i := 0; i < 10; i++ {
		// 使用limiter.Wait(ctx)等待，直到可以访问资源（访问文件/执行数据库查询等）
		if err := limiter.Wait(context.Background()); err != nil {
			log.Fatalf("Failed to wait for rate limiter: %v", err)
		}
		// 访问资源
		fmt.Printf("Accessing resource at %v\n", time.Now())
	}
}

// Demo_CurrentHandleWithHTTP 模拟HTTP请求限流
func Demo_CurrentHandleWithHTTP() {
	// 每秒最多处理 1 个请求，允许突发 2 个请求
	limiter := rate.NewLimiter(1, 2)
	http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		if limiter.Allow() {
			fmt.Fprintln(w, "Request allowed")
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	})
	fmt.Println("Server started at :8080")
	_ = http.ListenAndServe(":8080", nil)
}

// Demo_CurrentHandleWorks 模拟并发处理任务
func Demo_CurrentHandleWorks() {
	var wg sync.WaitGroup
	numWorkers := 5 // 模拟5个并发请求
	// 构造限流器：每10s向桶中新增一个令牌，桶里最多能存放2个令牌 => 每10s能处理一个任务，一定时间内最多能处理2个任务（令牌有剩余）
	var l = rate.NewLimiter(rate.Every(time.Second*10), 2)
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, l)
	}

	wg.Wait()
}

func worker(id int, wg *sync.WaitGroup, limiter *rate.Limiter) {
	defer wg.Done()
	if limiter.Allow() {
		fmt.Printf("Worker %d processed at %s\n", id, time.Now().Format("15:04:05.000"))
	} else {
		fmt.Printf("Worker %d rejected at %s\n", id, time.Now().Format("15:04:05.000"))
	}
}
