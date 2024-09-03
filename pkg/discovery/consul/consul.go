package consul

import (
	"context"
	"errors"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/lapkomo2018/goTwitterAuthService/pkg/discovery"
	"strconv"
	"strings"
)

// Registry defines a Consul-based service discovery registry
type Registry struct {
	client *consul.Client
}

// NewRegistry creates a new Consul-based service
// registry instance
func NewRegistry(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client: client}, nil
}

// Register creates a service record in the registry
func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in the format <host>:<port>, example: localhost:8080")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:    instanceID,
		Name:  serviceName,
		Port:  port,
		Check: &consul.AgentServiceCheck{CheckID: instanceID, TTL: "5s"},
	})
}

// Deregister removes a service record from the registry
func (r *Registry) Deregister(ctx context.Context, instanceID, _ string) error {
	return r.client.Agent().ServiceDeregister(instanceID)
}

// ServiceAddresses returns the list of addresses of active instances of the service
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	services, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(services) == 0 {
		return nil, discovery.ErrNotFound
	}

	addrs := make([]string, 0, len(services))
	for _, service := range services {
		addrs = append(addrs, fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port))
	}
	return addrs, nil
}

// ReportHealthyState is a push mechanism for
// reporting healthy state to the registry
func (r *Registry) ReportHealthyState(instanceID, _ string) error {
	return r.client.Agent().PassTTL(instanceID, "")
}
