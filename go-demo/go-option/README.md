# [Go] Option选项设计模式 — — 编程方式基础入门

# 1 介绍
在 Go 开发中，我们经常遇到需要处理多参数配置的场景。传统方法存在诸多痛点，例如：

问题一：参数过多且顺序敏感，导致我们调用时难以阅读，不知道每个参数背后对应的含义
```go
// 参数过多且顺序敏感
func NewServer(addr string, port int, timeout time.Duration, maxConns int, tls bool) {
    // ...
}

// 调用时难以阅读
srv := NewServer(":8080", 3306, 10*time.Second, 100, true)
```
问题二：新增改动，需要修改所有调用点
```go
// 新增参数需修改所有调用点
func NewServer(..., enableLog bool) // 新增参数破坏现有代码
```

这时就可以使用Go自带的Option方式编程。

Go Option主要有以下优势：
- ✅ 自描述性：命名选项明确参数含义
- ✅ 安全扩展：新增选项不影响现有调用
- ✅ 默认值处理：自动应用合理默认配置
- ✅ 参数验证：可在选项函数中实现验证逻辑
- ✅ 顺序无关：任意顺序传递选项参数

# 2 基础入门
## 2.1 定义类结构以及Option

```go
// 定义配置结构体
type Config struct {
    Timeout time.Duration
    MaxConn int
    TLS     bool
}

// 定义选项函数类型
type Option func(*Config)

// 步骤3：实现构造函数
func NewConfig(opts ...Option) *Config {
    // 设置默认值
    cfg := &Config{
        Timeout: 10 * time.Second,
        MaxConn: 100,
        TLS:     false,
    }
    
    // 应用所有选项
    for _, opt := range opts {
    	//因为opt本身就是func，所以这里相当于调用函数，入参为cfg struct
        opt(cfg)
    }
    return cfg
}
```

## 2.2 定义选项函数Withxx
```go
// 带参数的选项
func WithTimeout(t time.Duration) Option {
    return func(c *Config) {
        c.Timeout = t
    }
}

// 无参数的选项（开关功能）
func WithTLS() Option {
    return func(c *Config) {
        c.TLS = true
    }
}

// 带验证的选项
func WithMaxConn(n int) Option {
    return func(c *Config) {
        if n > 0 {
            c.MaxConn = n
        } // 否则保持默认值
    }
}
```
## 2.3 使用

```go
// 只使用默认值
defaultCfg := NewConfig()

// 覆盖部分默认值
customCfg := NewConfig(
    WithTimeout(30*time.Second),
    WithMaxConn(200),
)

// 启用特定功能
secureCfg := NewConfig(
    WithTLS(),
    WithTimeout(15*time.Second),
)
```
## 全部代码

```go
package main

import (
	"fmt"
	"time"
)

// 定义配置结构体
type Config struct {
	Timeout time.Duration
	MaxConn int
	TLS     bool
}

// 定义选项函数类型
type Option func(*Config)

// 实现构造函数
func NewConfig(opts ...Option) *Config {
	// 设置默认值
	cfg := &Config{
		Timeout: 10 * time.Second,
		MaxConn: 100,
		TLS:     false,
	}

	// 应用所有选项
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// 带参数的选项
func WithTimeout(t time.Duration) Option {
	return func(c *Config) {
		c.Timeout = t
	}
}

// 无参数的选项（开关功能）
func WithTLS() Option {
	return func(c *Config) {
		c.TLS = true
	}
}

// 带验证的选项
func WithMaxConn(n int) Option {
	return func(c *Config) {
		if n > 0 {
			c.MaxConn = n
		} // 否则保持默认值
	}
}

func main() {
	// 启用特定功能
	secureCfg := NewConfig(
		WithTLS(),
		WithTimeout(15*time.Second),
	)
	fmt.Println(secureCfg)
}
```

# 3. 实战使用
> Go Option方式可以用在数据库配置、HTTP服务配置、客户端连接配置、日志系统配置等。这里以HTTP服务配置为例。

```go
package main

import (
	"fmt"
	"time"
)

type Option func(*ServerConfig)

type ServerConfig struct {
	Addr        string
	ReadTimeout time.Duration
	IdleTimeout time.Duration
	EnableCORS  bool
}

func NewServer(addr string, opts ...Option) *ServerConfig {
	cfg := &ServerConfig{
		Addr:        addr,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 30 * time.Second,
	}

	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// 组合选项：同时设置多个相关参数
func WithTimeouts(read, idle time.Duration) Option {
	return func(s *ServerConfig) {
		s.ReadTimeout = read
		s.IdleTimeout = idle
	}
}

func EnableCORS() Option {
	return func(s *ServerConfig) {
		s.EnableCORS = true
	}
}

func main() {
	// 使用示例
	server := NewServer(
		":8080",
		WithTimeouts(10*time.Second, 60*time.Second),
		EnableCORS(),
	)
	fmt.Println(server)
}
```

# 4. 进阶（Option与链式调用结合）
> 在 Go 中我们可以结合 Option 模式和链式调用可以创建高度可读、灵活的 API，实现优雅编程。这种方式尤其适用于复杂对象的配置。
> 核心思想：
> 1. Option 模式：使用函数闭包封装配置逻辑
> 2. 链式调用：每个配置方法返回对象本身，支持连续调用

```go
package main

import "fmt"

// 目标配置对象
type Server struct {
	host    string
	port    int
	timeout int // 秒
	tls     bool
}

// Option 函数类型：接收 *Server 的闭包
type Option func(*Server)

// 链式包装器（关键结构）
type ServerBuilder struct {
	options []Option
}

// 创建 Builder 实例
func NewBuilder() *ServerBuilder {
	return &ServerBuilder{}
}

// 链式方法：添加配置选项
func (b *ServerBuilder) WithHost(host string) *ServerBuilder {
	b.options = append(b.options, func(s *Server) {
		s.host = host
	})
	return b
}

func (b *ServerBuilder) WithPort(port int) *ServerBuilder {
	b.options = append(b.options, func(s *Server) {
		s.port = port
	})
	return b
}

func (b *ServerBuilder) WithTimeout(timeout int) *ServerBuilder {
	b.options = append(b.options, func(s *Server) {
		s.timeout = timeout
	})
	return b
}

func (b *ServerBuilder) WithTLS(tls bool) *ServerBuilder {
	b.options = append(b.options, func(s *Server) {
		s.tls = tls
	})
	return b
}

func (b *ServerBuilder) Build() *Server {
	// 设置默认值
	s := &Server{
		host:    "localhost",
		port:    8080,
		timeout: 30,
	}

	// 应用所有配置函数
	for _, option := range b.options {
		option(s)
	}
	return s
}

func main() {
	server := NewBuilder().
		WithHost("api.ziyi.com").
		WithPort(443).
		WithTimeout(60).
		WithTLS(true).
		Build()

	fmt.Printf("%+v\n", server)
	// 输出：&{host:api.example.com port:443 timeout:60 tls:true}
}
```

# 总结
①概念：
> Option 模式是 Go 语言中处理复杂配置的优雅解决方案，它通过：
> - 功能选项（Functional Options）实现灵活的配置扩展
> - 默认值机制减少调用方负担
> - 命名参数提高代码可读性
> - 零成本扩展支持未来需求变化

②使用场景：
> 1. 对象需要 5个以上 配置参数时
> 2. 超过 3个可选 配置项时
> 3. 配置可能有 合理默认值 时
> 4. 需要 高频扩展 配置的场景
> 5. 开源库/框架中需要提供 友好API 时

实践tips：可以从简单的配置对象开始，当可选参数超过3个或发现构造函数参数过多时，可考虑重构为Option模式。