package user

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/auth/gateway/user/grpc"
)

type Gateway struct {
	grpc *grpc.Gateway
}

func New(httpAddr, grpcAddr string) *Gateway {
	return &Gateway{
		grpc: grpc.New(grpcAddr),
	}
}

func (g *Gateway) CreateUser(ctx context.Context, userID, role uint, username, email string) error {
	return g.grpc.CreateUser(ctx, userID, role, username, email)
}
