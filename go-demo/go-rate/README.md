# 【go类库分享】go rate标准限速库
> 在我们实际生产开发过程中，不免存在一些高并发的场景，此时我们就需要对请求进行限流，避免过高的QPS，影响我们服务器，导致服务出现波动或不可用。
> - 对于go开发人员而言，golang标准库就已经为我们提供了现成的类库`golang.org/x/time/rate`，我们只需直接调用即可。下面将为大家介绍该类库的详细用法。
## 介绍
> 限流（Rate Limiting）是控制对某些资源访问频率的一种技术手段。在高并发的服务中，限流机制可以有效防止资源过载、服务崩溃，保障系统的稳定性和可用性。Golang 官方标准库 golang.org/x/time/rate 提供了一个高效且易用的限流器（Rate Limiter），可以帮助开发者方便地实现限流功能。

原理：golang官方的类库采用了令牌桶算法进行限流。
- 令牌桶算法（Token Bucket Algorithm）：是一种常用的限流算法，它通过在固定时间间隔内向“桶”中添加“令牌”，请求在处理前需要从桶中获取令牌。如果桶中有足够的令牌，请求被处理；否则，请求被拒绝或等待。
- 速率（Rate）：速率定义了令牌添加的速度，即每秒向桶中添加多少令牌。
- 容量（Burst）：容量定义了桶的大小，即桶中最多可以存储多少令牌。它决定了在一段时间内允许的最大突发请求数。

## 安装

```go
go get golang.org/x/time/rate
```
## API介绍
### rate.NewLimiter：创建限流器
> 创建限流器后，可以通过 Allow、Reserve、Wait 等方法请求许可

```go
package main

import (
    "fmt"
    "golang.org/x/time/rate"
    "time"
)

func main() {
    // 每秒生成3个令牌，桶的容量为10个令牌
    // 相当于每秒最多能处理三个请求，顺时并发最多能处理10个请求
    limiter := rate.NewLimiter(3, 10)

    fmt.Println("Limiter created with rate 3 tokens per second and burst size of 10")
}
```

### limiter.Allow()：请求是否被允许/限流
> Allow 方法立即返回一个布尔值，指示请求是否被允许

```go
if limiter.Allow() {
    fmt.Println("Request allowed")
} else {
    fmt.Println("Request denied")
}
```

### limiter.Reserve()：返回值包含了许可时间和是否可用的信息
> Reserve 方法返回一个 Reservation 对象，包含了许可时间和是否可用的信息
```go
reservation := limiter.Reserve()
if reservation.OK() {
    fmt.Println("Request reserved, delay:", reservation.Delay())
} else {
    fmt.Println("Request cannot be reserved")
}
```

### limiter.Wait(ctx)：阻塞当前协程，直到允许请求或上下文取消
> Wait 方法阻塞当前协程，直到允许请求或上下文取消

```go
ctx := context.Background()
if err := limiter.Wait(ctx); err == nil {
    fmt.Println("Request allowed after wait")
} else {
    fmt.Println("Request denied:", err)
}
```

## 实战使用
### 并发处理任务
> 一段时间内，限制服务器处理任务数
```go

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
```
### HTTP服务请求限流
> 当我们通过jmeter/脚本/连续刷新页面并发请求测试时，发现在1s内，服务器最多只能处理2个请求，其余请求都不会被处理，会返回Too Many Requests。
```go
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
```


### 模拟等待耗时任务处理
> 可用于限制资源（文件）的访问，避免资源（文件）被频繁访问造成性能问题。
```go
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
```



参考文章：<br/>
https://cloud.tencent.com/developer/article/2429254<br/>
https://www.cnblogs.com/gnivor/p/10623028.html

