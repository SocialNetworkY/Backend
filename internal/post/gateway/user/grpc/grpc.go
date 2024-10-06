package grpc

import (
	"context"
	"time"

	"github.com/lapkomo2018/goTwitterServices/pkg/constant"
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

func (g *Gateway) GetUserRole(ctx context.Context, auth string, userID uint) (uint, error) {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewUserServiceClient(conn)

	resp, err := client.GetUserRole(grpcutil.PutMetadata(ctx, constant.GRPCAuthorizationMetadata, auth), &gen.GetUserRoleRequest{
		UserId: uint64(userID),
	})
	if err != nil {
		return 0, err
	}

	return uint(resp.GetRole()), nil
}

func (g *Gateway) IsUserBan(ctx context.Context, auth string, userID uint) (bool, string, time.Time, error) {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return false, "", time.Time{}, err
	}
	defer conn.Close()
	client := gen.NewUserServiceClient(conn)

	resp, err := client.IsUserBan(grpcutil.PutMetadata(ctx, constant.GRPCAuthorizationMetadata, auth), &gen.IsUserBanRequest{
		UserId: uint64(userID),
	})
	if err != nil {
		return false, "", time.Time{}, err
	}

	expiredAt, err := time.Parse(time.RFC3339, resp.GetExpiredAt())
	if err != nil {
		return false, "", time.Time{}, err
	}

	return resp.GetBanned(), resp.Reason, expiredAt, nil
}
