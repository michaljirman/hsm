// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// MpcSignerClient is the client API for MpcSigner service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MpcSignerClient interface {
	Ready(ctx context.Context, in *ReadyRequest, opts ...grpc.CallOption) (*ReadyResponse, error)
	Shutdown(ctx context.Context, in *ShutdownRequest, opts ...grpc.CallOption) (*ShutdownResponse, error)
	Signature(ctx context.Context, in *ShutdownRequest, opts ...grpc.CallOption) (*SignatureResponse, error)
	Test(ctx context.Context, opts ...grpc.CallOption) (MpcSigner_TestClient, error)
}

type mpcSignerClient struct {
	cc grpc.ClientConnInterface
}

func NewMpcSignerClient(cc grpc.ClientConnInterface) MpcSignerClient {
	return &mpcSignerClient{cc}
}

func (c *mpcSignerClient) Ready(ctx context.Context, in *ReadyRequest, opts ...grpc.CallOption) (*ReadyResponse, error) {
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, "/protobuf.MpcSigner/Ready", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mpcSignerClient) Shutdown(ctx context.Context, in *ShutdownRequest, opts ...grpc.CallOption) (*ShutdownResponse, error) {
	out := new(ShutdownResponse)
	err := c.cc.Invoke(ctx, "/protobuf.MpcSigner/Shutdown", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mpcSignerClient) Signature(ctx context.Context, in *ShutdownRequest, opts ...grpc.CallOption) (*SignatureResponse, error) {
	out := new(SignatureResponse)
	err := c.cc.Invoke(ctx, "/protobuf.MpcSigner/Signature", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mpcSignerClient) Test(ctx context.Context, opts ...grpc.CallOption) (MpcSigner_TestClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MpcSigner_serviceDesc.Streams[0], "/protobuf.MpcSigner/Test", opts...)
	if err != nil {
		return nil, err
	}
	x := &mpcSignerTestClient{stream}
	return x, nil
}

type MpcSigner_TestClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type mpcSignerTestClient struct {
	grpc.ClientStream
}

func (x *mpcSignerTestClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *mpcSignerTestClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MpcSignerServer is the server API for MpcSigner service.
// All implementations must embed UnimplementedMpcSignerServer
// for forward compatibility
type MpcSignerServer interface {
	Ready(context.Context, *ReadyRequest) (*ReadyResponse, error)
	Shutdown(context.Context, *ShutdownRequest) (*ShutdownResponse, error)
	Signature(context.Context, *ShutdownRequest) (*SignatureResponse, error)
	Test(MpcSigner_TestServer) error
	mustEmbedUnimplementedMpcSignerServer()
}

// UnimplementedMpcSignerServer must be embedded to have forward compatible implementations.
type UnimplementedMpcSignerServer struct {
}

func (UnimplementedMpcSignerServer) Ready(context.Context, *ReadyRequest) (*ReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ready not implemented")
}
func (UnimplementedMpcSignerServer) Shutdown(context.Context, *ShutdownRequest) (*ShutdownResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shutdown not implemented")
}
func (UnimplementedMpcSignerServer) Signature(context.Context, *ShutdownRequest) (*SignatureResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signature not implemented")
}
func (UnimplementedMpcSignerServer) Test(MpcSigner_TestServer) error {
	return status.Errorf(codes.Unimplemented, "method Test not implemented")
}
func (UnimplementedMpcSignerServer) mustEmbedUnimplementedMpcSignerServer() {}

// UnsafeMpcSignerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MpcSignerServer will
// result in compilation errors.
type UnsafeMpcSignerServer interface {
	mustEmbedUnimplementedMpcSignerServer()
}

func RegisterMpcSignerServer(s grpc.ServiceRegistrar, srv MpcSignerServer) {
	s.RegisterService(&_MpcSigner_serviceDesc, srv)
}

func _MpcSigner_Ready_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MpcSignerServer).Ready(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.MpcSigner/Ready",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MpcSignerServer).Ready(ctx, req.(*ReadyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MpcSigner_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShutdownRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MpcSignerServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.MpcSigner/Shutdown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MpcSignerServer).Shutdown(ctx, req.(*ShutdownRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MpcSigner_Signature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShutdownRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MpcSignerServer).Signature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.MpcSigner/Signature",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MpcSignerServer).Signature(ctx, req.(*ShutdownRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MpcSigner_Test_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MpcSignerServer).Test(&mpcSignerTestServer{stream})
}

type MpcSigner_TestServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type mpcSignerTestServer struct {
	grpc.ServerStream
}

func (x *mpcSignerTestServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *mpcSignerTestServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _MpcSigner_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.MpcSigner",
	HandlerType: (*MpcSignerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ready",
			Handler:    _MpcSigner_Ready_Handler,
		},
		{
			MethodName: "Shutdown",
			Handler:    _MpcSigner_Shutdown_Handler,
		},
		{
			MethodName: "Signature",
			Handler:    _MpcSigner_Signature_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Test",
			Handler:       _MpcSigner_Test_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "service.proto",
}
