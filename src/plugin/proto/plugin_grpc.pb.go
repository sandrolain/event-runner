// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: src/plugin/proto/plugin.proto

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

const (
	PluginService_Status_FullMethodName   = "/eventrunner.plugin.v1.PluginService/Status"
	PluginService_Command_FullMethodName  = "/eventrunner.plugin.v1.PluginService/Command"
	PluginService_Shutdown_FullMethodName = "/eventrunner.plugin.v1.PluginService/Shutdown"
	PluginService_Output_FullMethodName   = "/eventrunner.plugin.v1.PluginService/Output"
	PluginService_Input_FullMethodName    = "/eventrunner.plugin.v1.PluginService/Input"
)

// PluginServiceClient is the client API for PluginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PluginServiceClient interface {
	Status(ctx context.Context, in *StatusReq, opts ...grpc.CallOption) (*StatusRes, error)
	Command(ctx context.Context, in *CommandReq, opts ...grpc.CallOption) (*CommandRes, error)
	Shutdown(ctx context.Context, in *ShutdownReq, opts ...grpc.CallOption) (*ShutdownRes, error)
	Output(ctx context.Context, in *OutputReq, opts ...grpc.CallOption) (*OutputRes, error)
	Input(ctx context.Context, in *InputReq, opts ...grpc.CallOption) (PluginService_InputClient, error)
}

type pluginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginServiceClient(cc grpc.ClientConnInterface) PluginServiceClient {
	return &pluginServiceClient{cc}
}

func (c *pluginServiceClient) Status(ctx context.Context, in *StatusReq, opts ...grpc.CallOption) (*StatusRes, error) {
	out := new(StatusRes)
	err := c.cc.Invoke(ctx, PluginService_Status_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginServiceClient) Command(ctx context.Context, in *CommandReq, opts ...grpc.CallOption) (*CommandRes, error) {
	out := new(CommandRes)
	err := c.cc.Invoke(ctx, PluginService_Command_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginServiceClient) Shutdown(ctx context.Context, in *ShutdownReq, opts ...grpc.CallOption) (*ShutdownRes, error) {
	out := new(ShutdownRes)
	err := c.cc.Invoke(ctx, PluginService_Shutdown_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginServiceClient) Output(ctx context.Context, in *OutputReq, opts ...grpc.CallOption) (*OutputRes, error) {
	out := new(OutputRes)
	err := c.cc.Invoke(ctx, PluginService_Output_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginServiceClient) Input(ctx context.Context, in *InputReq, opts ...grpc.CallOption) (PluginService_InputClient, error) {
	stream, err := c.cc.NewStream(ctx, &PluginService_ServiceDesc.Streams[0], PluginService_Input_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &pluginServiceInputClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PluginService_InputClient interface {
	Recv() (*InputRes, error)
	grpc.ClientStream
}

type pluginServiceInputClient struct {
	grpc.ClientStream
}

func (x *pluginServiceInputClient) Recv() (*InputRes, error) {
	m := new(InputRes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PluginServiceServer is the server API for PluginService service.
// All implementations must embed UnimplementedPluginServiceServer
// for forward compatibility
type PluginServiceServer interface {
	Status(context.Context, *StatusReq) (*StatusRes, error)
	Command(context.Context, *CommandReq) (*CommandRes, error)
	Shutdown(context.Context, *ShutdownReq) (*ShutdownRes, error)
	Output(context.Context, *OutputReq) (*OutputRes, error)
	Input(*InputReq, PluginService_InputServer) error
	mustEmbedUnimplementedPluginServiceServer()
}

// UnimplementedPluginServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPluginServiceServer struct {
}

func (UnimplementedPluginServiceServer) Status(context.Context, *StatusReq) (*StatusRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedPluginServiceServer) Command(context.Context, *CommandReq) (*CommandRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Command not implemented")
}
func (UnimplementedPluginServiceServer) Shutdown(context.Context, *ShutdownReq) (*ShutdownRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shutdown not implemented")
}
func (UnimplementedPluginServiceServer) Output(context.Context, *OutputReq) (*OutputRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Output not implemented")
}
func (UnimplementedPluginServiceServer) Input(*InputReq, PluginService_InputServer) error {
	return status.Errorf(codes.Unimplemented, "method Input not implemented")
}
func (UnimplementedPluginServiceServer) mustEmbedUnimplementedPluginServiceServer() {}

// UnsafePluginServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PluginServiceServer will
// result in compilation errors.
type UnsafePluginServiceServer interface {
	mustEmbedUnimplementedPluginServiceServer()
}

func RegisterPluginServiceServer(s grpc.ServiceRegistrar, srv PluginServiceServer) {
	s.RegisterService(&PluginService_ServiceDesc, srv)
}

func _PluginService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginService_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServiceServer).Status(ctx, req.(*StatusReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginService_Command_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServiceServer).Command(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginService_Command_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServiceServer).Command(ctx, req.(*CommandReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginService_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShutdownReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServiceServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginService_Shutdown_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServiceServer).Shutdown(ctx, req.(*ShutdownReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginService_Output_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OutputReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServiceServer).Output(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginService_Output_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServiceServer).Output(ctx, req.(*OutputReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginService_Input_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(InputReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PluginServiceServer).Input(m, &pluginServiceInputServer{stream})
}

type PluginService_InputServer interface {
	Send(*InputRes) error
	grpc.ServerStream
}

type pluginServiceInputServer struct {
	grpc.ServerStream
}

func (x *pluginServiceInputServer) Send(m *InputRes) error {
	return x.ServerStream.SendMsg(m)
}

// PluginService_ServiceDesc is the grpc.ServiceDesc for PluginService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PluginService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "eventrunner.plugin.v1.PluginService",
	HandlerType: (*PluginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _PluginService_Status_Handler,
		},
		{
			MethodName: "Command",
			Handler:    _PluginService_Command_Handler,
		},
		{
			MethodName: "Shutdown",
			Handler:    _PluginService_Shutdown_Handler,
		},
		{
			MethodName: "Output",
			Handler:    _PluginService_Output_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Input",
			Handler:       _PluginService_Input_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "src/plugin/proto/plugin.proto",
}

const (
	AppService_Result_FullMethodName = "/eventrunner.plugin.v1.AppService/Result"
)

// AppServiceClient is the client API for AppService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AppServiceClient interface {
	Result(ctx context.Context, in *ResultReq, opts ...grpc.CallOption) (*ResultRes, error)
}

type appServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAppServiceClient(cc grpc.ClientConnInterface) AppServiceClient {
	return &appServiceClient{cc}
}

func (c *appServiceClient) Result(ctx context.Context, in *ResultReq, opts ...grpc.CallOption) (*ResultRes, error) {
	out := new(ResultRes)
	err := c.cc.Invoke(ctx, AppService_Result_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AppServiceServer is the server API for AppService service.
// All implementations must embed UnimplementedAppServiceServer
// for forward compatibility
type AppServiceServer interface {
	Result(context.Context, *ResultReq) (*ResultRes, error)
	mustEmbedUnimplementedAppServiceServer()
}

// UnimplementedAppServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAppServiceServer struct {
}

func (UnimplementedAppServiceServer) Result(context.Context, *ResultReq) (*ResultRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Result not implemented")
}
func (UnimplementedAppServiceServer) mustEmbedUnimplementedAppServiceServer() {}

// UnsafeAppServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppServiceServer will
// result in compilation errors.
type UnsafeAppServiceServer interface {
	mustEmbedUnimplementedAppServiceServer()
}

func RegisterAppServiceServer(s grpc.ServiceRegistrar, srv AppServiceServer) {
	s.RegisterService(&AppService_ServiceDesc, srv)
}

func _AppService_Result_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResultReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServiceServer).Result(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AppService_Result_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServiceServer).Result(ctx, req.(*ResultReq))
	}
	return interceptor(ctx, in, info, handler)
}

// AppService_ServiceDesc is the grpc.ServiceDesc for AppService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AppService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "eventrunner.plugin.v1.AppService",
	HandlerType: (*AppServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Result",
			Handler:    _AppService_Result_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/plugin/proto/plugin.proto",
}
