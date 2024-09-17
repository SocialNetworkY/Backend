package grpcutil

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/pkg/discovery"
	"google.golang.org/grpc"
	"math/rand"
)

// ServiceConnection returns a connection to a random service
// instance from the provided service name.
func ServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.ServiceAddresses(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(addrs[rand.Intn(len(addrs))])
}
