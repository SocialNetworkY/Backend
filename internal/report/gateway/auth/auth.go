package auth

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/report/gateway/auth/grpc"
)

type Gateway struct {
	grpc *grpc.Gateway
}

func New(httpAddr, grpcAddr string) *Gateway {
	return &Gateway{
		grpc: grpc.New(grpcAddr),
	}
}

func (g *Gateway) Authenticate(ctx context.Context, auth string) (uint, error) {
	return g.grpc.Authenticate(ctx, auth)
}
