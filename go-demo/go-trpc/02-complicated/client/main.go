package main

import (
	"context"
	"trpc.group/trpc-go/trpc-go/client"
	"ziyi.com/02-complicated/pb"
)

func main() {
	cli := pb.NewClassroomServiceClientProxy(client.WithTarget("ip://localhost:8088"))
	rsp, err := cli.GetInfo(context.TODO(), &pb.Request{RoomId: 1})
	if err != nil {
		panic(err)
	}
	println("【client】 receive ", rsp.Classroom.Name)
}
