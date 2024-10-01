package grpc

import (
	"context"
	"errors"

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

func (g *Gateway) CreateUser(ctx context.Context, auth string, userID, role uint, username, email string) error {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewUserServiceClient(conn)

	resp, err := client.CreateUser(grpcutil.PutAuth(ctx, auth), &gen.CreateUserRequest{
		UserId:   uint64(userID),
		Role:     uint64(role),
		Username: username,
		Email:    email,
	})
	if err != nil {
		return err
	}
	if resp.GetUserId() != uint64(userID) {
		return errors.New("user id mismatch")
	}

	return nil
}

func (g *Gateway) GetUserRole(ctx context.Context, auth string, userID uint) (uint, error) {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewUserServiceClient(conn)

	resp, err := client.GetUserRole(grpcutil.PutAuth(ctx, auth), &gen.GetUserRoleRequest{
		UserId: uint64(userID),
	})
	if err != nil {
		return 0, err
	}

	return uint(resp.GetRole()), nil
}
