package user

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/internal/post/gateway/user/grpc"
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

func (g *Gateway) GetUserRole(ctx context.Context, auth string, userID uint) (uint, error) {
	return g.grpc.GetUserRole(ctx, auth, userID)
}

func (g *Gateway) IsUserBan(ctx context.Context, auth string, userID uint) (banned bool, reason string, expiredAt time.Time, err error) {
	return g.grpc.IsUserBan(ctx, auth, userID)
}
