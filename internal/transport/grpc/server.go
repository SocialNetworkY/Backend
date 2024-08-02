package grpc

import (
	"fmt"
	"github.com/lapkomo2018/goTwitterAuthService/internal/transport/grpc/v1"
	grpcTest "github.com/lapkomo2018/goTwitterAuthService/pkg/grpc/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Server struct {
	port       int
	grpcServer *grpc.Server
}

func New(port int) *Server {
	log.Printf("Created grpc server with port: %d", port)
	grpcServ := grpc.NewServer()
	reflection.Register(grpcServ)

	return &Server{
		port:       port,
		grpcServer: grpcServ,
	}
}

func (s *Server) Init() *Server {
	grpcTest.RegisterTestServer(s.grpcServer, &v1.Handler{})

	return s
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	log.Printf("Grpc server listening at %v", lis.Addr())
	return s.grpcServer.Serve(lis)
}
