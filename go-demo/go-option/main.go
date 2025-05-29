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
