package main

import (
	"context"
	"trpc.group/trpc-go/trpc-go/client"
	"ziyi.com/03-keyword/pb"
)

func main() {
	cli := pb.NewKeywordServiceClientProxy(client.WithTarget("ip://localhost:8088"))
	req := &pb.Request{}
	req.ReqInfo = map[string]string{"name": "zhangsan"}
	// optional string reqCreateTime = 1; optional 标识字段可选，如果不设置，在序列化时不会包含该字段
	//createTime := "2022-01-01"
	//req.ReqCreateTime = &createTime
	// 序列化
	//data, err := proto.Marshal(req)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("=========marshal data....", string(data), "==========")
	rsp, err := cli.GetKeyword(context.TODO(), req)
	if err != nil {
		panic(err)
	}
	for k, v := range rsp.RspInfo {
		println("【client】 receive key: ", k, " value: ", v)
	}
}
