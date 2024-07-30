package server

import (
	"fmt"
	"github.com/lapkomo2018/goPetProjectServiceCollector/internal/transport/grpc/server/api"
	grpcTest "github.com/lapkomo2018/goPetProjectServiceCollector/pkg/grpc/test"
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
	grpcTest.RegisterTestServer(s.grpcServer, &api.Handler{})

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
