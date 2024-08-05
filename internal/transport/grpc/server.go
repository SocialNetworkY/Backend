package grpc

import (
	"fmt"
	"github.com/lapkomo2018/goTwitterAuthService/internal/transport/grpc/v1"
	grpcAuthService "github.com/lapkomo2018/goTwitterAuthService/pkg/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type (
	Config struct {
		Port int
	}

	Server struct {
		addr       string
		grpcServer *grpc.Server
	}
)

func New(config Config) *Server {
	log.Printf("Creating grpc server with port: %d", config.Port)
	grpcServ := grpc.NewServer()
	reflection.Register(grpcServ)

	return &Server{
		addr:       fmt.Sprintf(":%d", config.Port),
		grpcServer: grpcServ,
	}
}

func (s *Server) Init(authenticationService v1.AuthenticationService) *Server {
	handler := v1.New(authenticationService)
	grpcAuthService.RegisterAuthenticationServer(s.grpcServer, handler)
	return s
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	log.Printf("Grpc server listening at %v", lis.Addr())
	return s.grpcServer.Serve(lis)
}
