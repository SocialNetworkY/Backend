package discovery

import (
	"context"
	"errors"
)

type Registry interface {
	// Register a service with registry
	Register(ctx context.Context, instanceID, serviceName, hostPort string, tags []string) error
	// Deregister a service with registry
	Deregister(ctx context.Context, instanceID, serviceName string) error
	// ServiceAddresses returns the addresses of instances that provide the service
	ServiceAddresses(ctx context.Context, serviceName, tag string) ([]string, error)
	// ReportHealthyState reports the health status of the service
	ReportHealthyState(instanceID, serviceName string) error
}

// ErrNotFound is returned when the service is not found in the registry
var ErrNotFound = errors.New("no service addresses found")

// GenerateInstanceID generates a unique instance ID for the service
func GenerateInstanceID(serviceName, hostPort string) string {
	return serviceName + "@" + hostPort
}
