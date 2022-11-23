// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.14.0
// source: api/inspection.proto

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

// InspectServiceClient is the client API for InspectService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InspectServiceClient interface {
	Inspect(ctx context.Context, in *InspectRequest, opts ...grpc.CallOption) (*InspectResponse, error)
}

type inspectServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInspectServiceClient(cc grpc.ClientConnInterface) InspectServiceClient {
	return &inspectServiceClient{cc}
}

func (c *inspectServiceClient) Inspect(ctx context.Context, in *InspectRequest, opts ...grpc.CallOption) (*InspectResponse, error) {
	out := new(InspectResponse)
	err := c.cc.Invoke(ctx, "/api.InspectService/Inspect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InspectServiceServer is the server API for InspectService service.
// All implementations must embed UnimplementedInspectServiceServer
// for forward compatibility
type InspectServiceServer interface {
	Inspect(context.Context, *InspectRequest) (*InspectResponse, error)
	mustEmbedUnimplementedInspectServiceServer()
}

// UnimplementedInspectServiceServer must be embedded to have forward compatible implementations.
type UnimplementedInspectServiceServer struct {
}

func (UnimplementedInspectServiceServer) Inspect(context.Context, *InspectRequest) (*InspectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Inspect not implemented")
}
func (UnimplementedInspectServiceServer) mustEmbedUnimplementedInspectServiceServer() {}

// UnsafeInspectServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InspectServiceServer will
// result in compilation errors.
type UnsafeInspectServiceServer interface {
	mustEmbedUnimplementedInspectServiceServer()
}

func RegisterInspectServiceServer(s grpc.ServiceRegistrar, srv InspectServiceServer) {
	s.RegisterService(&InspectService_ServiceDesc, srv)
}

func _InspectService_Inspect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InspectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InspectServiceServer).Inspect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.InspectService/Inspect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InspectServiceServer).Inspect(ctx, req.(*InspectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InspectService_ServiceDesc is the grpc.ServiceDesc for InspectService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InspectService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.InspectService",
	HandlerType: (*InspectServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Inspect",
			Handler:    _InspectService_Inspect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/inspection.proto",
}
