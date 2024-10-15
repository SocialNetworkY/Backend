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