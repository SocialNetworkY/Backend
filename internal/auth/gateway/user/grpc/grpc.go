package grpc

import (
	"context"
	"errors"
	"github.com/SocialNetworkY/Backend/pkg/gen"
	"github.com/SocialNetworkY/Backend/pkg/grpcutil"
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

func (g *Gateway) CreateUser(ctx context.Context, userID, role uint, username, email string) error {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewUserServiceClient(conn)

	resp, err := client.CreateUser(ctx, &gen.CreateUserRequest{
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
