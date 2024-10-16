package post

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/user/gateway/post/grpc"
)

type Gateway struct {
	grpc *grpc.Gateway
}

func New(httpAddr, grpcAddr string) *Gateway {
	return &Gateway{
		grpc: grpc.New(grpcAddr),
	}
}

func (g *Gateway) DeleteUserPosts(ctx context.Context, id uint) error {
	return g.grpc.DeleteUserPosts(ctx, id)
}

func (g *Gateway) DeleteUserComments(ctx context.Context, id uint) error {
	return g.grpc.DeleteUserComments(ctx, id)
}

func (g *Gateway) DeleteUserLikes(ctx context.Context, id uint) error {
	return g.grpc.DeleteUserLikes(ctx, id)
}
