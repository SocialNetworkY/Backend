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

func (g *Gateway) DeleteUserPosts(ctx context.Context, id uint) error {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewPostServiceClient(conn)

	_, err = client.DeleteUserPosts(ctx, &gen.DeleteUserPostsRequest{UserId: uint64(id)})
	if err != nil {
		return err
	}

	return nil
}

func (g *Gateway) DeleteUserComments(ctx context.Context, id uint) error {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewPostServiceClient(conn)

	_, err = client.DeleteUserComments(ctx, &gen.DeleteUserCommentsRequest{UserId: uint64(id)})
	if err != nil {
		return err
	}

	return nil
}

func (g *Gateway) DeleteUserLikes(ctx context.Context, id uint) error {
	conn, err := grpcutil.ServiceConnection(g.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewPostServiceClient(conn)

	_, err = client.DeleteUserLikes(ctx, &gen.DeleteUserLikesRequest{UserId: uint64(id)})
	if err != nil {
		return err
	}

	return nil
}
