// Code generated by trpc-go/trpc-go-cmdline v2.6.5. DO NOT EDIT.
// source: classroom.proto

package pb

import (
	"context"
	"errors"
	"fmt"

	_ "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
	_ "trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/server"
)

// START ======================================= Server Service Definition ======================================= START

// ClassroomServiceService defines service.
type ClassroomServiceService interface {
	GetInfo(ctx context.Context, req *Request) (*Response, error)
}

func ClassroomServiceService_GetInfo_Handler(svr interface{}, ctx context.Context, f server.FilterFunc) (interface{}, error) {
	req := &Request{}
	filters, err := f(req)
	if err != nil {
		return nil, err
	}
	handleFunc := func(ctx context.Context, reqbody interface{}) (interface{}, error) {
		return svr.(ClassroomServiceService).GetInfo(ctx, reqbody.(*Request))
	}

	var rsp interface{}
	rsp, err = filters.Filter(ctx, req, handleFunc)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// ClassroomServiceServer_ServiceDesc descriptor for server.RegisterService.
var ClassroomServiceServer_ServiceDesc = server.ServiceDesc{
	ServiceName: "trpc.complicated.ClassroomService",
	HandlerType: ((*ClassroomServiceService)(nil)),
	Methods: []server.Method{
		{
			Name: "/trpc.complicated.ClassroomService/GetInfo",
			Func: ClassroomServiceService_GetInfo_Handler,
		},
	},
}

// RegisterClassroomServiceService registers service.
func RegisterClassroomServiceService(s server.Service, svr ClassroomServiceService) {
	if err := s.Register(&ClassroomServiceServer_ServiceDesc, svr); err != nil {
		panic(fmt.Sprintf("ClassroomService register error:%v", err))
	}
}

// START --------------------------------- Default Unimplemented Server Service --------------------------------- START

type UnimplementedClassroomService struct{}

func (s *UnimplementedClassroomService) GetInfo(ctx context.Context, req *Request) (*Response, error) {
	return nil, errors.New("rpc GetInfo of service ClassroomService is not implemented")
}

// END --------------------------------- Default Unimplemented Server Service --------------------------------- END

// END ======================================= Server Service Definition ======================================= END

// START ======================================= Client Service Definition ======================================= START

// ClassroomServiceClientProxy defines service client proxy
type ClassroomServiceClientProxy interface {
	GetInfo(ctx context.Context, req *Request, opts ...client.Option) (rsp *Response, err error)
}

type ClassroomServiceClientProxyImpl struct {
	client client.Client
	opts   []client.Option
}

var NewClassroomServiceClientProxy = func(opts ...client.Option) ClassroomServiceClientProxy {
	return &ClassroomServiceClientProxyImpl{client: client.DefaultClient, opts: opts}
}

func (c *ClassroomServiceClientProxyImpl) GetInfo(ctx context.Context, req *Request, opts ...client.Option) (*Response, error) {
	ctx, msg := codec.WithCloneMessage(ctx)
	defer codec.PutBackMessage(msg)
	msg.WithClientRPCName("/trpc.complicated.ClassroomService/GetInfo")
	msg.WithCalleeServiceName(ClassroomServiceServer_ServiceDesc.ServiceName)
	msg.WithCalleeApp("")
	msg.WithCalleeServer("")
	msg.WithCalleeService("ClassroomService")
	msg.WithCalleeMethod("GetInfo")
	msg.WithSerializationType(codec.SerializationTypePB)
	callopts := make([]client.Option, 0, len(c.opts)+len(opts))
	callopts = append(callopts, c.opts...)
	callopts = append(callopts, opts...)
	rsp := &Response{}
	if err := c.client.Invoke(ctx, req, rsp, callopts...); err != nil {
		return nil, err
	}
	return rsp, nil
}

// END ======================================= Client Service Definition ======================================= END