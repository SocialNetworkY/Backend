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

func (g *Gateway) DeletePostReports(ctx context.Context, postID uint) error {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewReportServiceClient(conn)

	_, err = client.DeletePostReports(ctx, &gen.DeletePostReportsRequest{
		PostId: uint64(postID),
	})
	return err
}
