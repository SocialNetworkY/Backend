package grpcutil

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/pkg/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"math/rand"
)

const (
	TagGRPC  = "grpc"
	MetaAuth = "authorization"
)

// ServiceConnection returns a connection to a random service
// instance from the provided service name.
func ServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.ServiceAddresses(ctx, serviceName, TagGRPC)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(addrs[rand.Intn(len(addrs))])
}

// PutAuth Put auth into context
func PutAuth(ctx context.Context, auth string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(MetaAuth, auth))
}
