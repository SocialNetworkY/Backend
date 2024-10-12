package grpc

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
	"time"

	"github.com/lapkomo2018/goTwitterServices/pkg/gen"
	"github.com/lapkomo2018/goTwitterServices/pkg/grpcutil"
)

// Gateway represents the gRPC gateway for the user service.
type Gateway struct {
	addr string
}

// New creates a new Gateway.
func New(addr string) *Gateway {
	return &Gateway{
		addr: addr,
	}
}

func (g *Gateway) UserInfo(ctx context.Context, userID uint) (*model.User, error) {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := gen.NewUserServiceClient(conn)

	resp, err := client.UserInfo(ctx, &gen.UserInfoRequest{
		UserId: uint64(userID),
	})
	if err != nil {
		return nil, err
	}

	banExpiredAt, _ := time.Parse(time.RFC3339, resp.GetBanExpiredAt())

	return &model.User{
		ID:           uint(resp.GetUserId()),
		Role:         uint(resp.GetRole()),
		Banned:       resp.GetBanned(),
		BanReason:    resp.GetBanReason(),
		BanExpiredAt: banExpiredAt,
	}, nil
}
