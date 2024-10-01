package user

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/gateway/user/grpc"
)

type Gateway struct {
	grpc *grpc.Gateway
}

func New(httpAddr, grpcAddr string) *Gateway {
	return &Gateway{
		grpc: grpc.New(grpcAddr),
	}
}

func (g *Gateway) CreateUser(ctx context.Context, auth string, userID, role uint, username, email string) error {
	return g.grpc.CreateUser(ctx, auth, userID, role, username, email)
}

func (g *Gateway) GetUserRole(ctx context.Context, auth string, userID uint) (uint, error) {
	return g.grpc.GetUserRole(ctx, auth, userID)
}
