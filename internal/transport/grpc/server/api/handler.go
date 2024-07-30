package api

import (
	"context"
	grpcTest "github.com/lapkomo2018/goPetProjectServiceCollector/pkg/grpc/test"
	"log"
)

type Handler struct {
	grpcTest.UnimplementedTestServer
}

func (s *Handler) Get(ctx context.Context, r *grpcTest.GetRequest) (*grpcTest.GetResponse, error) {
	log.Printf("Get message: %v", r.Message)
	return &grpcTest.GetResponse{Message: r.Message}, nil
}
