package user

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/gateway/user/grpc"
	"time"
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
	_, userRole, _, _, _, err := g.UserInfo(ctx, userID)
	return userRole, err
}

func (g *Gateway) UserInfo(ctx context.Context, userID uint) (uint, uint, bool, string, time.Time, error) {
	return g.grpc.UserInfo(ctx, userID)
}
