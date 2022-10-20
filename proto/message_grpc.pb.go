// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proto/message.proto

package proto

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

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	CombineMsg(ctx context.Context, in *MessageUploadCombine, opts ...grpc.CallOption) (*ConfirmationFromNameNode, error)
	ToDataNodeMsg(ctx context.Context, in *MessageUploadToDataNode, opts ...grpc.CallOption) (*ConfirmationFromDataNode, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) CombineMsg(ctx context.Context, in *MessageUploadCombine, opts ...grpc.CallOption) (*ConfirmationFromNameNode, error) {
	out := new(ConfirmationFromNameNode)
	err := c.cc.Invoke(ctx, "/grpc.MessageService/CombineMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) ToDataNodeMsg(ctx context.Context, in *MessageUploadToDataNode, opts ...grpc.CallOption) (*ConfirmationFromDataNode, error) {
	out := new(ConfirmationFromDataNode)
	err := c.cc.Invoke(ctx, "/grpc.MessageService/ToDataNodeMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations must embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	CombineMsg(context.Context, *MessageUploadCombine) (*ConfirmationFromNameNode, error)
	ToDataNodeMsg(context.Context, *MessageUploadToDataNode) (*ConfirmationFromDataNode, error)
	mustEmbedUnimplementedMessageServiceServer()
}

// UnimplementedMessageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) CombineMsg(context.Context, *MessageUploadCombine) (*ConfirmationFromNameNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CombineMsg not implemented")
}
func (UnimplementedMessageServiceServer) ToDataNodeMsg(context.Context, *MessageUploadToDataNode) (*ConfirmationFromDataNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToDataNodeMsg not implemented")
}
func (UnimplementedMessageServiceServer) mustEmbedUnimplementedMessageServiceServer() {}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_CombineMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageUploadCombine)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).CombineMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.MessageService/CombineMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).CombineMsg(ctx, req.(*MessageUploadCombine))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_ToDataNodeMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageUploadToDataNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).ToDataNodeMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.MessageService/ToDataNodeMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).ToDataNodeMsg(ctx, req.(*MessageUploadToDataNode))
	}
	return interceptor(ctx, in, info, handler)
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CombineMsg",
			Handler:    _MessageService_CombineMsg_Handler,
		},
		{
			MethodName: "ToDataNodeMsg",
			Handler:    _MessageService_ToDataNodeMsg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/message.proto",
}