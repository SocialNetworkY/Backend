package grpc

import "github.com/lapkomo2018/goTwitterServices/pkg/discovery"

// Gateway represents the gRPC gateway for the user service.
type Gateway struct {
	registry discovery.Registry
}

const UserServiceName = "user"

// New creates a new Gateway.
func New(r discovery.Registry) *Gateway {
	return &Gateway{registry: r}
}
