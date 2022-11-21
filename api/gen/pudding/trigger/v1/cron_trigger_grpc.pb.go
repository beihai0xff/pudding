// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// CronTriggerServiceClient is the client API for CronTriggerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CronTriggerServiceClient interface {
	// Sends a Ping
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingResponse, error)
	// UpdateStatus update cron trigger status
	UpdateStatus(ctx context.Context, in *UpdateStatusRequest, opts ...grpc.CallOption) (*UpdateStatusResponse, error)
}

type cronTriggerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCronTriggerServiceClient(cc grpc.ClientConnInterface) CronTriggerServiceClient {
	return &cronTriggerServiceClient{cc}
}

func (c *cronTriggerServiceClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/pudding.trigger.v1.CronTriggerService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cronTriggerServiceClient) UpdateStatus(ctx context.Context, in *UpdateStatusRequest, opts ...grpc.CallOption) (*UpdateStatusResponse, error) {
	out := new(UpdateStatusResponse)
	err := c.cc.Invoke(ctx, "/pudding.trigger.v1.CronTriggerService/UpdateStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CronTriggerServiceServer is the server API for CronTriggerService service.
// All implementations must embed UnimplementedCronTriggerServiceServer
// for forward compatibility
type CronTriggerServiceServer interface {
	// Sends a Ping
	Ping(context.Context, *emptypb.Empty) (*PingResponse, error)
	// UpdateStatus update cron trigger status
	UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateStatusResponse, error)
	mustEmbedUnimplementedCronTriggerServiceServer()
}

// UnimplementedCronTriggerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCronTriggerServiceServer struct {
}

func (UnimplementedCronTriggerServiceServer) Ping(context.Context, *emptypb.Empty) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedCronTriggerServiceServer) UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatus not implemented")
}
func (UnimplementedCronTriggerServiceServer) mustEmbedUnimplementedCronTriggerServiceServer() {}

// UnsafeCronTriggerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CronTriggerServiceServer will
// result in compilation errors.
type UnsafeCronTriggerServiceServer interface {
	mustEmbedUnimplementedCronTriggerServiceServer()
}

func RegisterCronTriggerServiceServer(s grpc.ServiceRegistrar, srv CronTriggerServiceServer) {
	s.RegisterService(&CronTriggerService_ServiceDesc, srv)
}

func _CronTriggerService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronTriggerServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pudding.trigger.v1.CronTriggerService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronTriggerServiceServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CronTriggerService_UpdateStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronTriggerServiceServer).UpdateStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pudding.trigger.v1.CronTriggerService/UpdateStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronTriggerServiceServer).UpdateStatus(ctx, req.(*UpdateStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CronTriggerService_ServiceDesc is the grpc.ServiceDesc for CronTriggerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CronTriggerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pudding.trigger.v1.CronTriggerService",
	HandlerType: (*CronTriggerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _CronTriggerService_Ping_Handler,
		},
		{
			MethodName: "UpdateStatus",
			Handler:    _CronTriggerService_UpdateStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pudding/trigger/v1/cron_trigger.proto",
}
