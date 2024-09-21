package main

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"trpc.group/trpc-go/trpc-go"
	"ziyi.com/05-validate-self/pb"
)

func main() {
	//设置服务配置文件路径
	trpc.ServerConfigPath = "/Users/ziyi/GolandProjects/ziyifast-code_instruction/go-demo/go-trpc/05-validate-self/server/trpc_go.yaml"
	server := trpc.NewServer()
	pb.RegisterUserServiceService(server, &UserService{})
	if err := server.Serve(); err != nil {
		panic(err)
	}
}

type UserService struct {
}

func (c *UserService) Handle(ctx context.Context, req *pb.User) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	log.Printf("pass.....")
	return nil, nil
}
