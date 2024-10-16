package grpc

import (
	"context"
	"github.com/SocialNetworkY/Backend/pkg/gen"
	"github.com/SocialNetworkY/Backend/pkg/grpcutil"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{
		addr: addr,
	}
}

func (g *Gateway) DeleteUserReports(ctx context.Context, userID uint) error {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewReportServiceClient(conn)

	_, err = client.DeleteUserReports(ctx, &gen.DeleteUserReportsRequest{
		UserId: uint64(userID),
	})
	return err
}
