package user

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/report/gateway/user/grpc"
	"github.com/SocialNetworkY/Backend/internal/report/model"
)

type Gateway struct {
	grpc *grpc.Gateway
}

func New(httpAddr, grpcAddr string) *Gateway {
	return &Gateway{
		grpc: grpc.New(grpcAddr),
	}
}

func (g *Gateway) UserInfo(ctx context.Context, userID uint) (*model.User, error) {
	return g.grpc.UserInfo(ctx, userID)
}
