// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/messages.proto

package messages

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
	AskAddress(ctx context.Context, in *Cmd, opts ...grpc.CallOption) (*Address, error)
	Command(ctx context.Context, in *Cmd, opts ...grpc.CallOption) (*VectorClock, error)
	GetSoldados(ctx context.Context, in *Info, opts ...grpc.CallOption) (*ValueInfo, error)
	Soldados(ctx context.Context, in *Info, opts ...grpc.CallOption) (*ValueInfo, error)
	RequestMerge(ctx context.Context, in *AskMerge, opts ...grpc.CallOption) (*Ack, error)
	Merge(ctx context.Context, in *AskMerge, opts ...grpc.CallOption) (*AllInfo, error)
	PropagateChanges(ctx context.Context, in *CommandsAndVectorClocks, opts ...grpc.CallOption) (*Ack, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) AskAddress(ctx context.Context, in *Cmd, opts ...grpc.CallOption) (*Address, error) {
	out := new(Address)
	err := c.cc.Invoke(ctx, "/messages.MessageService/AskAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) Command(ctx context.Context, in *Cmd, opts ...grpc.CallOption) (*VectorClock, error) {
	out := new(VectorClock)
	err := c.cc.Invoke(ctx, "/messages.MessageService/Command", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) GetSoldados(ctx context.Context, in *Info, opts ...grpc.CallOption) (*ValueInfo, error) {
	out := new(ValueInfo)
	err := c.cc.Invoke(ctx, "/messages.MessageService/GetSoldados", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) Soldados(ctx context.Context, in *Info, opts ...grpc.CallOption) (*ValueInfo, error) {
	out := new(ValueInfo)
	err := c.cc.Invoke(ctx, "/messages.MessageService/Soldados", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) RequestMerge(ctx context.Context, in *AskMerge, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/messages.MessageService/RequestMerge", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) Merge(ctx context.Context, in *AskMerge, opts ...grpc.CallOption) (*AllInfo, error) {
	out := new(AllInfo)
	err := c.cc.Invoke(ctx, "/messages.MessageService/Merge", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) PropagateChanges(ctx context.Context, in *CommandsAndVectorClocks, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/messages.MessageService/PropagateChanges", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations must embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	AskAddress(context.Context, *Cmd) (*Address, error)
	Command(context.Context, *Cmd) (*VectorClock, error)
	GetSoldados(context.Context, *Info) (*ValueInfo, error)
	Soldados(context.Context, *Info) (*ValueInfo, error)
	RequestMerge(context.Context, *AskMerge) (*Ack, error)
	Merge(context.Context, *AskMerge) (*AllInfo, error)
	PropagateChanges(context.Context, *CommandsAndVectorClocks) (*Ack, error)
	mustEmbedUnimplementedMessageServiceServer()
}

// UnimplementedMessageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) AskAddress(context.Context, *Cmd) (*Address, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AskAddress not implemented")
}
func (UnimplementedMessageServiceServer) Command(context.Context, *Cmd) (*VectorClock, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Command not implemented")
}
func (UnimplementedMessageServiceServer) GetSoldados(context.Context, *Info) (*ValueInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSoldados not implemented")
}
func (UnimplementedMessageServiceServer) Soldados(context.Context, *Info) (*ValueInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Soldados not implemented")
}
func (UnimplementedMessageServiceServer) RequestMerge(context.Context, *AskMerge) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestMerge not implemented")
}
func (UnimplementedMessageServiceServer) Merge(context.Context, *AskMerge) (*AllInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Merge not implemented")
}
func (UnimplementedMessageServiceServer) PropagateChanges(context.Context, *CommandsAndVectorClocks) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PropagateChanges not implemented")
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

func _MessageService_AskAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Cmd)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).AskAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.MessageService/AskAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).AskAddress(ctx, req.(*Cmd))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_Command_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Cmd)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).Command(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.MessageService/Command",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).Command(ctx, req.(*Cmd))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_GetSoldados_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Info)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).GetSoldados(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.MessageService/GetSoldados",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).GetSoldados(ctx, req.(*Info))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_Soldados_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Info)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).Soldados(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.MessageService/Soldados",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).Soldados(ctx, req.(*Info))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_RequestMerge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AskMerge)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).RequestMerge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.MessageService/RequestMerge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).RequestMerge(ctx, req.(*AskMerge))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_Merge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AskMerge)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).Merge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.MessageService/Merge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).Merge(ctx, req.(*AskMerge))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_PropagateChanges_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandsAndVectorClocks)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).PropagateChanges(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.MessageService/PropagateChanges",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).PropagateChanges(ctx, req.(*CommandsAndVectorClocks))
	}
	return interceptor(ctx, in, info, handler)
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "messages.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AskAddress",
			Handler:    _MessageService_AskAddress_Handler,
		},
		{
			MethodName: "Command",
			Handler:    _MessageService_Command_Handler,
		},
		{
			MethodName: "GetSoldados",
			Handler:    _MessageService_GetSoldados_Handler,
		},
		{
			MethodName: "Soldados",
			Handler:    _MessageService_Soldados_Handler,
		},
		{
			MethodName: "RequestMerge",
			Handler:    _MessageService_RequestMerge_Handler,
		},
		{
			MethodName: "Merge",
			Handler:    _MessageService_Merge_Handler,
		},
		{
			MethodName: "PropagateChanges",
			Handler:    _MessageService_PropagateChanges_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/messages.proto",
}