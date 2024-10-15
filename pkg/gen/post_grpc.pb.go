// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.2
// source: api/post.proto

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

// PostServiceClient is the client API for PostService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostServiceClient interface {
	DeleteUserPosts(ctx context.Context, in *DeleteUserPostsRequest, opts ...grpc.CallOption) (*DeleteUserPostsResponse, error)
}

type postServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPostServiceClient(cc grpc.ClientConnInterface) PostServiceClient {
	return &postServiceClient{cc}
}

func (c *postServiceClient) DeleteUserPosts(ctx context.Context, in *DeleteUserPostsRequest, opts ...grpc.CallOption) (*DeleteUserPostsResponse, error) {
	out := new(DeleteUserPostsResponse)
	err := c.cc.Invoke(ctx, "/PostService/DeleteUserPosts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostServiceServer is the server API for PostService service.
// All implementations must embed UnimplementedPostServiceServer
// for forward compatibility
type PostServiceServer interface {
	DeleteUserPosts(context.Context, *DeleteUserPostsRequest) (*DeleteUserPostsResponse, error)
	mustEmbedUnimplementedPostServiceServer()
}

// UnimplementedPostServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPostServiceServer struct {
}

func (UnimplementedPostServiceServer) DeleteUserPosts(context.Context, *DeleteUserPostsRequest) (*DeleteUserPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserPosts not implemented")
}
func (UnimplementedPostServiceServer) mustEmbedUnimplementedPostServiceServer() {}

// UnsafePostServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostServiceServer will
// result in compilation errors.
type UnsafePostServiceServer interface {
	mustEmbedUnimplementedPostServiceServer()
}

func RegisterPostServiceServer(s grpc.ServiceRegistrar, srv PostServiceServer) {
	s.RegisterService(&PostService_ServiceDesc, srv)
}

func _PostService_DeleteUserPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserPostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).DeleteUserPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PostService/DeleteUserPosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).DeleteUserPosts(ctx, req.(*DeleteUserPostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PostService_ServiceDesc is the grpc.ServiceDesc for PostService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "PostService",
	HandlerType: (*PostServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteUserPosts",
			Handler:    _PostService_DeleteUserPosts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/post.proto",
}