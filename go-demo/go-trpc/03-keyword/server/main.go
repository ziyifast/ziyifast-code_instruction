package main

import (
	"context"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
	"ziyi.com/03-keyword/pb"
)

func main() {
	trpc.ServerConfigPath = "/Users/ziyi/GolandProjects/ziyifast-code_instruction/go-demo/go-trpc/03-keyword/server/trpc_go.yaml"
	server := trpc.NewServer()
	pb.RegisterKeywordServiceService(server, &KeywordService{})
	if err := server.Serve(); err != nil {
		panic(err)
	}
}

type KeywordService struct {
}

func (k *KeywordService) GetKeyword(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Info("【server】 receive reqInfo：", req.ReqInfo, " reqTime ", req.ReqCreateTime)
	classroom := new(pb.Classroom)
	classroom.StudentIds = []int32{1, 2, 3}
	rsp := new(pb.Response)
	rsp.RspInfo = map[string]string{
		"data": classroom.String(),
	}
	log.Infof("%v", rsp.RspInfo["data"])
	return rsp, nil
}
