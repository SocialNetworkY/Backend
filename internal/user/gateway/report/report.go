package report

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/user/gateway/report/grpc"
)

type Gateway struct {
	grpc *grpc.Gateway
}

func New(httpAddr, grpcAddr string) *Gateway {
	return &Gateway{
		grpc: grpc.New(grpcAddr),
	}
}

func (g *Gateway) DeleteUserReports(ctx context.Context, userID uint) error {
	return g.grpc.DeleteUserReports(ctx, userID)
}
