// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.11
// source: client/api/daemon.proto

package api

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

// DaemonClient is the client API for Daemon service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DaemonClient interface {
	GetVersion(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Version, error)
	PingTest(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	AddProxyProfile(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error)
	GetProxyProfiles(ctx context.Context, in *Request, opts ...grpc.CallOption) (*GetProxyProfilesResponse, error)
	NotifyProxy(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyResponse, error)
	PullImage(ctx context.Context, in *ImageReference, opts ...grpc.CallOption) (*ImagePullResponse, error)
	SetOptimizer(ctx context.Context, in *OptimizeRequest, opts ...grpc.CallOption) (*OptimizeResponse, error)
	ReportTraces(ctx context.Context, in *ReportTracesRequest, opts ...grpc.CallOption) (*ReportTracesResponse, error)
}

type daemonClient struct {
	cc grpc.ClientConnInterface
}

func NewDaemonClient(cc grpc.ClientConnInterface) DaemonClient {
	return &daemonClient{cc}
}

func (c *daemonClient) GetVersion(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Version, error) {
	out := new(Version)
	err := c.cc.Invoke(ctx, "/api.Daemon/GetVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) PingTest(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/api.Daemon/PingTest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) AddProxyProfile(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, "/api.Daemon/AddProxyProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) GetProxyProfiles(ctx context.Context, in *Request, opts ...grpc.CallOption) (*GetProxyProfilesResponse, error) {
	out := new(GetProxyProfilesResponse)
	err := c.cc.Invoke(ctx, "/api.Daemon/GetProxyProfiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) NotifyProxy(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyResponse, error) {
	out := new(NotifyResponse)
	err := c.cc.Invoke(ctx, "/api.Daemon/NotifyProxy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) PullImage(ctx context.Context, in *ImageReference, opts ...grpc.CallOption) (*ImagePullResponse, error) {
	out := new(ImagePullResponse)
	err := c.cc.Invoke(ctx, "/api.Daemon/PullImage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) SetOptimizer(ctx context.Context, in *OptimizeRequest, opts ...grpc.CallOption) (*OptimizeResponse, error) {
	out := new(OptimizeResponse)
	err := c.cc.Invoke(ctx, "/api.Daemon/SetOptimizer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) ReportTraces(ctx context.Context, in *ReportTracesRequest, opts ...grpc.CallOption) (*ReportTracesResponse, error) {
	out := new(ReportTracesResponse)
	err := c.cc.Invoke(ctx, "/api.Daemon/ReportTraces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DaemonServer is the server API for Daemon service.
// All implementations must embed UnimplementedDaemonServer
// for forward compatibility
type DaemonServer interface {
	GetVersion(context.Context, *Request) (*Version, error)
	PingTest(context.Context, *PingRequest) (*PingResponse, error)
	AddProxyProfile(context.Context, *AuthRequest) (*AuthResponse, error)
	GetProxyProfiles(context.Context, *Request) (*GetProxyProfilesResponse, error)
	NotifyProxy(context.Context, *NotifyRequest) (*NotifyResponse, error)
	PullImage(context.Context, *ImageReference) (*ImagePullResponse, error)
	SetOptimizer(context.Context, *OptimizeRequest) (*OptimizeResponse, error)
	ReportTraces(context.Context, *ReportTracesRequest) (*ReportTracesResponse, error)
	mustEmbedUnimplementedDaemonServer()
}

// UnimplementedDaemonServer must be embedded to have forward compatible implementations.
type UnimplementedDaemonServer struct {
}

func (UnimplementedDaemonServer) GetVersion(context.Context, *Request) (*Version, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (UnimplementedDaemonServer) PingTest(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingTest not implemented")
}
func (UnimplementedDaemonServer) AddProxyProfile(context.Context, *AuthRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProxyProfile not implemented")
}
func (UnimplementedDaemonServer) GetProxyProfiles(context.Context, *Request) (*GetProxyProfilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProxyProfiles not implemented")
}
func (UnimplementedDaemonServer) NotifyProxy(context.Context, *NotifyRequest) (*NotifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyProxy not implemented")
}
func (UnimplementedDaemonServer) PullImage(context.Context, *ImageReference) (*ImagePullResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PullImage not implemented")
}
func (UnimplementedDaemonServer) SetOptimizer(context.Context, *OptimizeRequest) (*OptimizeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOptimizer not implemented")
}
func (UnimplementedDaemonServer) ReportTraces(context.Context, *ReportTracesRequest) (*ReportTracesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportTraces not implemented")
}
func (UnimplementedDaemonServer) mustEmbedUnimplementedDaemonServer() {}

// UnsafeDaemonServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DaemonServer will
// result in compilation errors.
type UnsafeDaemonServer interface {
	mustEmbedUnimplementedDaemonServer()
}

func RegisterDaemonServer(s grpc.ServiceRegistrar, srv DaemonServer) {
	s.RegisterService(&Daemon_ServiceDesc, srv)
}

func _Daemon_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Daemon/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).GetVersion(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_PingTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).PingTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Daemon/PingTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).PingTest(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_AddProxyProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).AddProxyProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Daemon/AddProxyProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).AddProxyProfile(ctx, req.(*AuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_GetProxyProfiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).GetProxyProfiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Daemon/GetProxyProfiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).GetProxyProfiles(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_NotifyProxy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).NotifyProxy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Daemon/NotifyProxy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).NotifyProxy(ctx, req.(*NotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_PullImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImageReference)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).PullImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Daemon/PullImage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).PullImage(ctx, req.(*ImageReference))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_SetOptimizer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OptimizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).SetOptimizer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Daemon/SetOptimizer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).SetOptimizer(ctx, req.(*OptimizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_ReportTraces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportTracesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).ReportTraces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Daemon/ReportTraces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).ReportTraces(ctx, req.(*ReportTracesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Daemon_ServiceDesc is the grpc.ServiceDesc for Daemon service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Daemon_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Daemon",
	HandlerType: (*DaemonServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVersion",
			Handler:    _Daemon_GetVersion_Handler,
		},
		{
			MethodName: "PingTest",
			Handler:    _Daemon_PingTest_Handler,
		},
		{
			MethodName: "AddProxyProfile",
			Handler:    _Daemon_AddProxyProfile_Handler,
		},
		{
			MethodName: "GetProxyProfiles",
			Handler:    _Daemon_GetProxyProfiles_Handler,
		},
		{
			MethodName: "NotifyProxy",
			Handler:    _Daemon_NotifyProxy_Handler,
		},
		{
			MethodName: "PullImage",
			Handler:    _Daemon_PullImage_Handler,
		},
		{
			MethodName: "SetOptimizer",
			Handler:    _Daemon_SetOptimizer_Handler,
		},
		{
			MethodName: "ReportTraces",
			Handler:    _Daemon_ReportTraces_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "client/api/daemon.proto",
}
