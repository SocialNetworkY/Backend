package grpc

import (
	"context"
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

func (g *Gateway) Authenticate(ctx context.Context, auth string) (uint, error) {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewAuthServiceClient(conn)

	resp, err := client.Authenticate(ctx, &gen.AuthenticateRequest{
		AuthToken: auth,
	})
	if err != nil {
		return 0, err
	}

	return uint(resp.GetUserId()), nil
}
