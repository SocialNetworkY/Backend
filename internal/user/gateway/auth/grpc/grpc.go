package grpc

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/pkg/discovery"
	"github.com/lapkomo2018/goTwitterServices/pkg/gen"
	"github.com/lapkomo2018/goTwitterServices/pkg/grpcutil"
)

// Gateway represents the gRPC gateway for the user service.
type Gateway struct {
	registry discovery.Registry
}

const (
	AuthServiceName = "AuthService"
)

// New creates a new Gateway.
func New(r discovery.Registry) *Gateway {
	return &Gateway{registry: r}
}

func (g *Gateway) CheckAuth(ctx context.Context, auth string) (uint, error) {
	conn, err := grpcutil.ServiceConnection(ctx, AuthServiceName, g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewAuthServiceClient(conn)

	resp, err := client.Authenticate(grpcutil.PutAuth(ctx, auth), &gen.AuthenticateRequest{})
	if err != nil {
		return 0, err
	}

	return uint(resp.UserId), nil
}

func (g *Gateway) UpdateUsernameEmail(ctx context.Context, auth string, id uint, username, email string) error {
	conn, err := grpcutil.ServiceConnection(ctx, AuthServiceName, g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewAuthServiceClient(conn)

	_, err = client.UpdateUsernameEmail(grpcutil.PutAuth(ctx, auth), &gen.UpdateUsernameEmailRequest{
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
	conn, err := grpcutil.ServiceConnection(ctx, AuthServiceName, g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewAuthServiceClient(conn)

	_, err = client.DeleteUser(grpcutil.PutAuth(ctx, auth), &gen.DeleteUserRequest{UserId: uint64(id)})
	if err != nil {
		return err
	}

	return nil
}
