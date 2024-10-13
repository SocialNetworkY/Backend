package grpc

import (
	"context"
	"fmt"
	"github.com/SocialNetworkY/Backend/internal/user/transport/grpc/v1"
	"github.com/SocialNetworkY/Backend/pkg/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

type Server struct {
	addr       string
	grpcServer *grpc.Server
}

func New(port int) *Server {
	log.Printf("Creating grpc server with port: %d", port)
	grpcServ := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor()),
	)
	reflection.Register(grpcServ)

	return &Server{
		addr:       fmt.Sprintf(":%d", port),
		grpcServer: grpcServ,
	}
}

func (s *Server) Init(us v1.UserService, ag v1.AuthGateway) *Server {
	handler := v1.New(us, ag)
	gen.RegisterUserServiceServer(s.grpcServer, handler)
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

// UnaryServerInterceptor for logging
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		h, err := handler(ctx, req)
		end := time.Now()

		log.Printf("Request - Method:%s\tDuration:%s\tError:%v\n",
			info.FullMethod,
			end.Sub(start),
			err)

		return h, err
	}
}
