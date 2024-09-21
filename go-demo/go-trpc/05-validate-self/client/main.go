package main

import (
	"context"
	"trpc.group/trpc-go/trpc-go/client"
	"ziyi.com/05-validate-self/pb"
)

func main() {
	cli := pb.NewUserServiceClientProxy(client.WithTarget("ip://localhost:8088"))
	req := new(pb.User)
	req.Name = "akajerry"
	_, err := cli.Handle(context.Background(), req)
	if err != nil {
		panic(err)
	}
}
