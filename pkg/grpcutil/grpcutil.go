package grpcutil

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"math/rand"
)

type Registry interface {
	// ServiceAddresses returns the addresses of instances that provide the service
	ServiceAddresses(ctx context.Context, serviceName, tag string) ([]string, error)
}

const (
	TagGRPC = "grpc"
)

// ServiceConnectionWithRegistry returns a connection to a random service
// instance from the provided service name.
func ServiceConnectionWithRegistry(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.ServiceAddresses(ctx, serviceName, TagGRPC)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(addrs[rand.Intn(len(addrs))], grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// ServiceConnection returns a connection to a random service instance from the provided service name.
func ServiceConnection(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// PutMetadata returns a new context with the provided key-value pair.
func PutMetadata(ctx context.Context, key, value string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(key, value))
}
