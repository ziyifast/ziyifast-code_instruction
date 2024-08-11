package main

import (
	"context"
	pb2 "github.com/go-demo/go-trpc/01-simple/pb"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	//设置server配置文件路径，默认在./trpc_go.yaml
	trpc.ServerConfigPath = "E:\\Go\\GoPro\\src\\go_code\\ziyifast-code_instruction\\go-demo\\go-trpc\\01-simple\\server\\trpc_go.yaml"
	s := trpc.NewServer()
	pb2.RegisterGreeterService(s, &Greeter{})
	if err := s.Serve(); err != nil {
		log.Error(err)
	}
}

type Greeter struct{}

// Hello API
// 1. 接受client请求并打印
// 2. 拼接Hello后作为响应返回给client
func (g Greeter) Hello(ctx context.Context, req *pb2.HelloRequest) (*pb2.HelloReply, error) {
	log.Infof("got hello request: %s", req.Msg)
	return &pb2.HelloReply{Msg: "Hello " + req.Msg + "!"}, nil
}
