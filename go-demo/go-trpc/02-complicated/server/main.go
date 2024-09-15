package main

import (
	"context"
	"fmt"
	"trpc.group/trpc-go/tnet/log"
	"trpc.group/trpc-go/trpc-go"
	"ziyi.com/02-complicated/pb"
)

func main() {
	trpc.ServerConfigPath = "/Users/ziyi/GolandProjects/ziyifast-code_instruction/go-demo/go-trpc/02-complicated/server/trpc_go.yaml"
	server := trpc.NewServer()
	pb.RegisterClassroomServiceService(server, &ClassRoomService{})
	if err := server.Serve(); err != nil {
		panic(err)
	}
}

type ClassRoomService struct {
}

func (c *ClassRoomService) GetInfo(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Info("【server】receive ", request.RoomId, " info...")
	if request.RoomId != 1 {
		return nil, fmt.Errorf("the classroom does not exist")
	}
	rsp := new(pb.Response)
	room := &pb.Classroom{
		Name:    "grade7_21",
		Address: "北京市 朝阳区 大屯路 ",
		Students: []*pb.Student{
			&pb.Student{
				Name: "小明",
				Age:  18,
			},
		},
	}
	rsp.Classroom = room
	return rsp, nil
}
