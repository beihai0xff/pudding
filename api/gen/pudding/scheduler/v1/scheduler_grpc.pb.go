// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package scheduler

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SchedulerServiceClient is the client API for SchedulerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SchedulerServiceClient interface {
	// Sends a Ping
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingResponse, error)
	// Sends a Delay Message
	SendDelayMessage(ctx context.Context, in *SendDelayMessageRequest, opts ...grpc.CallOption) (*SendDelayMessageResponse, error)
}

type schedulerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSchedulerServiceClient(cc grpc.ClientConnInterface) SchedulerServiceClient {
	return &schedulerServiceClient{cc}
}

func (c *schedulerServiceClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/pudding.scheduler.v1.SchedulerService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerServiceClient) SendDelayMessage(ctx context.Context, in *SendDelayMessageRequest, opts ...grpc.CallOption) (*SendDelayMessageResponse, error) {
	out := new(SendDelayMessageResponse)
	err := c.cc.Invoke(ctx, "/pudding.scheduler.v1.SchedulerService/SendDelayMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SchedulerServiceServer is the server API for SchedulerService service.
// All implementations must embed UnimplementedSchedulerServiceServer
// for forward compatibility
type SchedulerServiceServer interface {
	// Sends a Ping
	Ping(context.Context, *emptypb.Empty) (*PingResponse, error)
	// Sends a Delay Message
	SendDelayMessage(context.Context, *SendDelayMessageRequest) (*SendDelayMessageResponse, error)
	mustEmbedUnimplementedSchedulerServiceServer()
}

// UnimplementedSchedulerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSchedulerServiceServer struct {
}

func (UnimplementedSchedulerServiceServer) Ping(context.Context, *emptypb.Empty) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedSchedulerServiceServer) SendDelayMessage(context.Context, *SendDelayMessageRequest) (*SendDelayMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendDelayMessage not implemented")
}
func (UnimplementedSchedulerServiceServer) mustEmbedUnimplementedSchedulerServiceServer() {}

// UnsafeSchedulerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SchedulerServiceServer will
// result in compilation errors.
type UnsafeSchedulerServiceServer interface {
	mustEmbedUnimplementedSchedulerServiceServer()
}

func RegisterSchedulerServiceServer(s grpc.ServiceRegistrar, srv SchedulerServiceServer) {
	s.RegisterService(&SchedulerService_ServiceDesc, srv)
}

func _SchedulerService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pudding.scheduler.v1.SchedulerService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServiceServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchedulerService_SendDelayMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendDelayMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServiceServer).SendDelayMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pudding.scheduler.v1.SchedulerService/SendDelayMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServiceServer).SendDelayMessage(ctx, req.(*SendDelayMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SchedulerService_ServiceDesc is the grpc.ServiceDesc for SchedulerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SchedulerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pudding.scheduler.v1.SchedulerService",
	HandlerType: (*SchedulerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _SchedulerService_Ping_Handler,
		},
		{
			MethodName: "SendDelayMessage",
			Handler:    _SchedulerService_SendDelayMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pudding/scheduler/v1/scheduler.proto",
}
