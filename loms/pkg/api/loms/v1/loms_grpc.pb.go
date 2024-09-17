// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: api/loms/v1/loms.proto

package loms

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Loms_OrderCreate_FullMethodName = "/Loms/OrderCreate"
	Loms_OrderPay_FullMethodName    = "/Loms/OrderPay"
	Loms_OrderCancel_FullMethodName = "/Loms/OrderCancel"
	Loms_OrderInfo_FullMethodName   = "/Loms/OrderInfo"
	Loms_StocksInfo_FullMethodName  = "/Loms/StocksInfo"
)

// LomsClient is the client API for Loms service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LomsClient interface {
	OrderCreate(ctx context.Context, in *OrderCreateRequest, opts ...grpc.CallOption) (*OrderCreateResponse, error)
	OrderPay(ctx context.Context, in *OrderPayRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	OrderCancel(ctx context.Context, in *OrderCancelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	OrderInfo(ctx context.Context, in *OrderInfoRequest, opts ...grpc.CallOption) (*OrderInfoResponse, error)
	StocksInfo(ctx context.Context, in *StocksInfoRequest, opts ...grpc.CallOption) (*StocksInfoResponse, error)
}

type lomsClient struct {
	cc grpc.ClientConnInterface
}

func NewLomsClient(cc grpc.ClientConnInterface) LomsClient {
	return &lomsClient{cc}
}

func (c *lomsClient) OrderCreate(ctx context.Context, in *OrderCreateRequest, opts ...grpc.CallOption) (*OrderCreateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderCreateResponse)
	err := c.cc.Invoke(ctx, Loms_OrderCreate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsClient) OrderPay(ctx context.Context, in *OrderPayRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Loms_OrderPay_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsClient) OrderCancel(ctx context.Context, in *OrderCancelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Loms_OrderCancel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsClient) OrderInfo(ctx context.Context, in *OrderInfoRequest, opts ...grpc.CallOption) (*OrderInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderInfoResponse)
	err := c.cc.Invoke(ctx, Loms_OrderInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsClient) StocksInfo(ctx context.Context, in *StocksInfoRequest, opts ...grpc.CallOption) (*StocksInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StocksInfoResponse)
	err := c.cc.Invoke(ctx, Loms_StocksInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LomsServer is the server API for Loms service.
// All implementations must embed UnimplementedLomsServer
// for forward compatibility.
type LomsServer interface {
	OrderCreate(context.Context, *OrderCreateRequest) (*OrderCreateResponse, error)
	OrderPay(context.Context, *OrderPayRequest) (*emptypb.Empty, error)
	OrderCancel(context.Context, *OrderCancelRequest) (*emptypb.Empty, error)
	OrderInfo(context.Context, *OrderInfoRequest) (*OrderInfoResponse, error)
	StocksInfo(context.Context, *StocksInfoRequest) (*StocksInfoResponse, error)
	mustEmbedUnimplementedLomsServer()
}

// UnimplementedLomsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLomsServer struct{}

func (UnimplementedLomsServer) OrderCreate(context.Context, *OrderCreateRequest) (*OrderCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderCreate not implemented")
}
func (UnimplementedLomsServer) OrderPay(context.Context, *OrderPayRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderPay not implemented")
}
func (UnimplementedLomsServer) OrderCancel(context.Context, *OrderCancelRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderCancel not implemented")
}
func (UnimplementedLomsServer) OrderInfo(context.Context, *OrderInfoRequest) (*OrderInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderInfo not implemented")
}
func (UnimplementedLomsServer) StocksInfo(context.Context, *StocksInfoRequest) (*StocksInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StocksInfo not implemented")
}
func (UnimplementedLomsServer) mustEmbedUnimplementedLomsServer() {}
func (UnimplementedLomsServer) testEmbeddedByValue()              {}

// UnsafeLomsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LomsServer will
// result in compilation errors.
type UnsafeLomsServer interface {
	mustEmbedUnimplementedLomsServer()
}

func RegisterLomsServer(s grpc.ServiceRegistrar, srv LomsServer) {
	// If the following call pancis, it indicates UnimplementedLomsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Loms_ServiceDesc, srv)
}

func _Loms_OrderCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).OrderCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Loms_OrderCreate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).OrderCreate(ctx, req.(*OrderCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Loms_OrderPay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderPayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).OrderPay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Loms_OrderPay_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).OrderPay(ctx, req.(*OrderPayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Loms_OrderCancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderCancelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).OrderCancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Loms_OrderCancel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).OrderCancel(ctx, req.(*OrderCancelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Loms_OrderInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).OrderInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Loms_OrderInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).OrderInfo(ctx, req.(*OrderInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Loms_StocksInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StocksInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).StocksInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Loms_StocksInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).StocksInfo(ctx, req.(*StocksInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Loms_ServiceDesc is the grpc.ServiceDesc for Loms service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Loms_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Loms",
	HandlerType: (*LomsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OrderCreate",
			Handler:    _Loms_OrderCreate_Handler,
		},
		{
			MethodName: "OrderPay",
			Handler:    _Loms_OrderPay_Handler,
		},
		{
			MethodName: "OrderCancel",
			Handler:    _Loms_OrderCancel_Handler,
		},
		{
			MethodName: "OrderInfo",
			Handler:    _Loms_OrderInfo_Handler,
		},
		{
			MethodName: "StocksInfo",
			Handler:    _Loms_StocksInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/loms/v1/loms.proto",
}
