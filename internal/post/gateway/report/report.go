package report

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/post/gateway/report/grpc"
)

type Gateway struct {
	grpc *grpc.Gateway
}

func New(httpAddr, grpcAddr string) *Gateway {
	return &Gateway{
		grpc: grpc.New(grpcAddr),
	}
}

func (g *Gateway) DeletePostReports(ctx context.Context, postID uint) error {
	return g.grpc.DeletePostReports(ctx, postID)
}
