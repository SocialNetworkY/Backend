// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.2
// source: api/report.proto

package gen

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

// ReportServiceClient is the client API for ReportService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReportServiceClient interface {
	DeleteUserReports(ctx context.Context, in *DeleteUserReportsRequest, opts ...grpc.CallOption) (*DeleteUserReportsResponse, error)
	DeletePostReports(ctx context.Context, in *DeletePostReportsRequest, opts ...grpc.CallOption) (*DeletePostReportsResponse, error)
}

type reportServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReportServiceClient(cc grpc.ClientConnInterface) ReportServiceClient {
	return &reportServiceClient{cc}
}

func (c *reportServiceClient) DeleteUserReports(ctx context.Context, in *DeleteUserReportsRequest, opts ...grpc.CallOption) (*DeleteUserReportsResponse, error) {
	out := new(DeleteUserReportsResponse)
	err := c.cc.Invoke(ctx, "/ReportService/DeleteUserReports", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reportServiceClient) DeletePostReports(ctx context.Context, in *DeletePostReportsRequest, opts ...grpc.CallOption) (*DeletePostReportsResponse, error) {
	out := new(DeletePostReportsResponse)
	err := c.cc.Invoke(ctx, "/ReportService/DeletePostReports", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReportServiceServer is the server API for ReportService service.
// All implementations must embed UnimplementedReportServiceServer
// for forward compatibility
type ReportServiceServer interface {
	DeleteUserReports(context.Context, *DeleteUserReportsRequest) (*DeleteUserReportsResponse, error)
	DeletePostReports(context.Context, *DeletePostReportsRequest) (*DeletePostReportsResponse, error)
	mustEmbedUnimplementedReportServiceServer()
}

// UnimplementedReportServiceServer must be embedded to have forward compatible implementations.
type UnimplementedReportServiceServer struct {
}

func (UnimplementedReportServiceServer) DeleteUserReports(context.Context, *DeleteUserReportsRequest) (*DeleteUserReportsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserReports not implemented")
}
func (UnimplementedReportServiceServer) DeletePostReports(context.Context, *DeletePostReportsRequest) (*DeletePostReportsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePostReports not implemented")
}
func (UnimplementedReportServiceServer) mustEmbedUnimplementedReportServiceServer() {}

// UnsafeReportServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReportServiceServer will
// result in compilation errors.
type UnsafeReportServiceServer interface {
	mustEmbedUnimplementedReportServiceServer()
}

func RegisterReportServiceServer(s grpc.ServiceRegistrar, srv ReportServiceServer) {
	s.RegisterService(&ReportService_ServiceDesc, srv)
}

func _ReportService_DeleteUserReports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserReportsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServiceServer).DeleteUserReports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ReportService/DeleteUserReports",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServiceServer).DeleteUserReports(ctx, req.(*DeleteUserReportsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReportService_DeletePostReports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePostReportsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServiceServer).DeletePostReports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ReportService/DeletePostReports",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServiceServer).DeletePostReports(ctx, req.(*DeletePostReportsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReportService_ServiceDesc is the grpc.ServiceDesc for ReportService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReportService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ReportService",
	HandlerType: (*ReportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteUserReports",
			Handler:    _ReportService_DeleteUserReports_Handler,
		},
		{
			MethodName: "DeletePostReports",
			Handler:    _ReportService_DeletePostReports_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/report.proto",
}