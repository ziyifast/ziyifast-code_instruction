package main

import (
	"context"
	"log"
	"trpc.group/trpc-go/trpc-go"
	"ziyi.com/04-validate/pb"
)

func main() {
	//设置服务配置文件路径
	trpc.ServerConfigPath = "/Users/ziyi/GolandProjects/ziyifast-code_instruction/go-demo/go-trpc/04-validate/server/trpc_go.yaml"
	server := trpc.NewServer()
	pb.RegisterUserServiceService(server, &UserService{})
	if err := server.Serve(); err != nil {
		panic(err)
	}
}

type UserService struct {
}

func (s *UserService) GetUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	log.Printf("pass.....")
	return req, nil
}
