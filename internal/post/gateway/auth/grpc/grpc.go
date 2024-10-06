package grpc

import (
	"context"
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

func (g *Gateway) Authenticate(ctx context.Context, auth string) (uint, error) {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewAuthServiceClient(conn)

	resp, err := client.Authenticate(grpcutil.PutMetadata(ctx, constant.GRPCAuthorizationMetadata, auth), &gen.AuthenticateRequest{})
	if err != nil {
		return 0, err
	}

	return uint(resp.UserId), nil
}

func (g *Gateway) UpdateUsernameEmail(ctx context.Context, auth string, id uint, username, email string) error {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewAuthServiceClient(conn)

	_, err = client.UpdateUsernameEmail(grpcutil.PutMetadata(ctx, constant.GRPCAuthorizationMetadata, auth), &gen.UpdateUsernameEmailRequest{
		UserId:   uint64(id),
		Username: username,
		Email:    email,
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *Gateway) DeleteUser(ctx context.Context, auth string, id uint) error {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewAuthServiceClient(conn)

	_, err = client.DeleteUser(grpcutil.PutMetadata(ctx, constant.GRPCAuthorizationMetadata, auth), &gen.DeleteUserRequest{UserId: uint64(id)})
	if err != nil {
		return err
	}

	return nil
}
