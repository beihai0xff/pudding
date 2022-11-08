// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: scheduler.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PuddingClient is the client API for Pudding service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PuddingClient interface {
	// Sends a Ping
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReply, error)
	// Sends a Delay Message
	SendDelayMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageReply, error)
}

type puddingClient struct {
	cc grpc.ClientConnInterface
}

func NewPuddingClient(cc grpc.ClientConnInterface) PuddingClient {
	return &puddingClient{cc}
}

func (c *puddingClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReply, error) {
	out := new(PingReply)
	err := c.cc.Invoke(ctx, "/Pudding/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *puddingClient) SendDelayMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageReply, error) {
	out := new(SendMessageReply)
	err := c.cc.Invoke(ctx, "/Pudding/SendDelayMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PuddingServer is the server API for Pudding service.
// All implementations must embed UnimplementedPuddingServer
// for forward compatibility
type PuddingServer interface {
	// Sends a Ping
	Ping(context.Context, *PingRequest) (*PingReply, error)
	// Sends a Delay Message
	SendDelayMessage(context.Context, *SendMessageRequest) (*SendMessageReply, error)
	mustEmbedUnimplementedPuddingServer()
}

// UnimplementedPuddingServer must be embedded to have forward compatible implementations.
type UnimplementedPuddingServer struct {
}

func (UnimplementedPuddingServer) Ping(context.Context, *PingRequest) (*PingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedPuddingServer) SendDelayMessage(context.Context, *SendMessageRequest) (*SendMessageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendDelayMessage not implemented")
}
func (UnimplementedPuddingServer) mustEmbedUnimplementedPuddingServer() {}

// UnsafePuddingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PuddingServer will
// result in compilation errors.
type UnsafePuddingServer interface {
	mustEmbedUnimplementedPuddingServer()
}

func RegisterPuddingServer(s grpc.ServiceRegistrar, srv PuddingServer) {
	s.RegisterService(&Pudding_ServiceDesc, srv)
}

func _Pudding_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PuddingServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Pudding/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PuddingServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pudding_SendDelayMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PuddingServer).SendDelayMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Pudding/SendDelayMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PuddingServer).SendDelayMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Pudding_ServiceDesc is the grpc.ServiceDesc for Pudding service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Pudding_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Pudding",
	HandlerType: (*PuddingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Pudding_Ping_Handler,
		},
		{
			MethodName: "SendDelayMessage",
			Handler:    _Pudding_SendDelayMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "scheduler.proto",
}