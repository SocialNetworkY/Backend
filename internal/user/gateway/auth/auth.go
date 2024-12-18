package auth

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/user/gateway/auth/grpc"
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

func (g *Gateway) UpdateUsernameEmail(ctx context.Context, id uint, username, email string) error {
	return g.grpc.UpdateUsernameEmail(ctx, id, username, email)
}

func (g *Gateway) DeleteUser(ctx context.Context, id uint) error {
	return g.grpc.DeleteUser(ctx, id)
}
