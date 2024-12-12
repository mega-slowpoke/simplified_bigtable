// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.24.4
// source: tablet.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	TabletService_Read_FullMethodName   = "/bigtable.TabletService/Read"
	TabletService_Write_FullMethodName  = "/bigtable.TabletService/Write"
	TabletService_Delete_FullMethodName = "/bigtable.TabletService/Delete"
	TabletService_Scan_FullMethodName   = "/bigtable.TabletService/Scan"
)

// TabletServiceClient is the client API for TabletService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Service for client interactions with user tablet servers
type TabletServiceClient interface {
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (*ScanResponse, error)
}

type tabletServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTabletServiceClient(cc grpc.ClientConnInterface) TabletServiceClient {
	return &tabletServiceClient{cc}
}

func (c *tabletServiceClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, TabletService_Read_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tabletServiceClient) Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WriteResponse)
	err := c.cc.Invoke(ctx, TabletService_Write_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tabletServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, TabletService_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tabletServiceClient) Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (*ScanResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ScanResponse)
	err := c.cc.Invoke(ctx, TabletService_Scan_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TabletServiceServer is the server API for TabletService service.
// All implementations must embed UnimplementedTabletServiceServer
// for forward compatibility.
//
// Service for client interactions with user tablet servers
type TabletServiceServer interface {
	Read(context.Context, *ReadRequest) (*ReadResponse, error)
	Write(context.Context, *WriteRequest) (*WriteResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	Scan(context.Context, *ScanRequest) (*ScanResponse, error)
	mustEmbedUnimplementedTabletServiceServer()
}

// UnimplementedTabletServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTabletServiceServer struct{}

func (UnimplementedTabletServiceServer) Read(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedTabletServiceServer) Write(context.Context, *WriteRequest) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedTabletServiceServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedTabletServiceServer) Scan(context.Context, *ScanRequest) (*ScanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Scan not implemented")
}
func (UnimplementedTabletServiceServer) mustEmbedUnimplementedTabletServiceServer() {}
func (UnimplementedTabletServiceServer) testEmbeddedByValue()                       {}

// UnsafeTabletServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TabletServiceServer will
// result in compilation errors.
type UnsafeTabletServiceServer interface {
	mustEmbedUnimplementedTabletServiceServer()
}

func RegisterTabletServiceServer(s grpc.ServiceRegistrar, srv TabletServiceServer) {
	// If the following call pancis, it indicates UnimplementedTabletServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TabletService_ServiceDesc, srv)
}

func _TabletService_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TabletServiceServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TabletService_Read_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TabletServiceServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TabletService_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TabletServiceServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TabletService_Write_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TabletServiceServer).Write(ctx, req.(*WriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TabletService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TabletServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TabletService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TabletServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TabletService_Scan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TabletServiceServer).Scan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TabletService_Scan_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TabletServiceServer).Scan(ctx, req.(*ScanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TabletService_ServiceDesc is the grpc.ServiceDesc for TabletService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TabletService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bigtable.TabletService",
	HandlerType: (*TabletServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Read",
			Handler:    _TabletService_Read_Handler,
		},
		{
			MethodName: "Write",
			Handler:    _TabletService_Write_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _TabletService_Delete_Handler,
		},
		{
			MethodName: "Scan",
			Handler:    _TabletService_Scan_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tablet.proto",
}
