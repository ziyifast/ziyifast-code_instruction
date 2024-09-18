package main

import (
	"context"
	"trpc.group/trpc-go/trpc-go/client"
	"ziyi.com/04-validate/pb"
)

func main() {
	cli := pb.NewUserServiceClientProxy(client.WithTarget("ip://localhost:8088"))
	req := new(pb.User)
	req.Email = "123@qq.com"
	req.Phone = "aaaa" //校验不通过
	//req.Phone = "18173827109" //校验通过
	req.Uid = 10001
	_, err := cli.GetUser(context.Background(), req)
	if err != nil {
		panic(err)
	}
}
