// Code generated by trpc-go/trpc-go-cmdline v2.6.5. DO NOT EDIT.
// source: user.proto

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

// UserServiceService defines service.
type UserServiceService interface {
	GetUser(ctx context.Context, req *User) (*User, error)
}

func UserServiceService_GetUser_Handler(svr interface{}, ctx context.Context, f server.FilterFunc) (interface{}, error) {
	req := &User{}
	filters, err := f(req)
	if err != nil {
		return nil, err
	}
	handleFunc := func(ctx context.Context, reqbody interface{}) (interface{}, error) {
		return svr.(UserServiceService).GetUser(ctx, reqbody.(*User))
	}

	var rsp interface{}
	rsp, err = filters.Filter(ctx, req, handleFunc)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// UserServiceServer_ServiceDesc descriptor for server.RegisterService.
var UserServiceServer_ServiceDesc = server.ServiceDesc{
	ServiceName: "validate.UserService",
	HandlerType: ((*UserServiceService)(nil)),
	Methods: []server.Method{
		{
			Name: "/validate.UserService/GetUser",
			Func: UserServiceService_GetUser_Handler,
		},
	},
}

// RegisterUserServiceService registers service.
func RegisterUserServiceService(s server.Service, svr UserServiceService) {
	if err := s.Register(&UserServiceServer_ServiceDesc, svr); err != nil {
		panic(fmt.Sprintf("UserService register error:%v", err))
	}
}

// START --------------------------------- Default Unimplemented Server Service --------------------------------- START

type UnimplementedUserService struct{}

func (s *UnimplementedUserService) GetUser(ctx context.Context, req *User) (*User, error) {
	return nil, errors.New("rpc GetUser of service UserService is not implemented")
}

// END --------------------------------- Default Unimplemented Server Service --------------------------------- END

// END ======================================= Server Service Definition ======================================= END

// START ======================================= Client Service Definition ======================================= START

// UserServiceClientProxy defines service client proxy
type UserServiceClientProxy interface {
	GetUser(ctx context.Context, req *User, opts ...client.Option) (rsp *User, err error)
}

type UserServiceClientProxyImpl struct {
	client client.Client
	opts   []client.Option
}

var NewUserServiceClientProxy = func(opts ...client.Option) UserServiceClientProxy {
	return &UserServiceClientProxyImpl{client: client.DefaultClient, opts: opts}
}

func (c *UserServiceClientProxyImpl) GetUser(ctx context.Context, req *User, opts ...client.Option) (*User, error) {
	ctx, msg := codec.WithCloneMessage(ctx)
	defer codec.PutBackMessage(msg)
	msg.WithClientRPCName("/validate.UserService/GetUser")
	msg.WithCalleeServiceName(UserServiceServer_ServiceDesc.ServiceName)
	msg.WithCalleeApp("")
	msg.WithCalleeServer("")
	msg.WithCalleeService("UserService")
	msg.WithCalleeMethod("GetUser")
	msg.WithSerializationType(codec.SerializationTypePB)
	callopts := make([]client.Option, 0, len(c.opts)+len(opts))
	callopts = append(callopts, c.opts...)
	callopts = append(callopts, opts...)
	rsp := &User{}
	if err := c.client.Invoke(ctx, req, rsp, callopts...); err != nil {
		return nil, err
	}
	return rsp, nil
}

// END ======================================= Client Service Definition ======================================= END
