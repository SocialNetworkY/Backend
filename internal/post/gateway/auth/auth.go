package auth

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/internal/post/gateway/auth/grpc"
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

func (g *Gateway) UpdateUsernameEmail(ctx context.Context, auth string, id uint, username, email string) error {
	return g.grpc.UpdateUsernameEmail(ctx, auth, id, username, email)
}

func (g *Gateway) DeleteUser(ctx context.Context, auth string, id uint) error {
	return g.grpc.DeleteUser(ctx, auth, id)
}
