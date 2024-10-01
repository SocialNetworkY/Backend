package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	authService "github.com/lapkomo2018/goTwitterServices/internal/user/gateway/auth/grpc"
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"github.com/lapkomo2018/goTwitterServices/internal/user/repository/mysql"
	"github.com/lapkomo2018/goTwitterServices/internal/user/service"
	"github.com/lapkomo2018/goTwitterServices/internal/user/transport/grpc"
	"github.com/lapkomo2018/goTwitterServices/internal/user/transport/rest"

	"github.com/lapkomo2018/goTwitterServices/config"
	"github.com/lapkomo2018/goTwitterServices/pkg/discovery"
	"github.com/lapkomo2018/goTwitterServices/pkg/discovery/consul"
)

type Config struct {
	Service    config.Service
	RestServer rest.Config
	GrpcServer grpc.Config
}

const (
	TagRest = "rest"
	TagGRPC = "grpc"
)

var (
	cfg *Config
	env *model.Env
)

func init() {
	var err error

	env, err = model.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err = config.LoadConfig[Config]()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	registry, err := consul.NewRegistry(env.DiscoveryAddr)
	if err != nil {
		log.Fatalf("could not create registry: %v", err)
	}
	ctx := context.Background()

	// Register REST service
	restInstanceID := discovery.GenerateInstanceID(cfg.Service.Name, env.ExternalRestPort)
	if err := registry.Register(ctx, restInstanceID, cfg.Service.Name, env.ExternalRestPort, []string{TagRest}); err != nil {
		log.Fatalf("could not register rest service: %v", err)
	}
	defer registry.Deregister(ctx, restInstanceID, cfg.Service.Name)

	// Register gRPC service
	grpcInstanceID := discovery.GenerateInstanceID(cfg.Service.Name, env.ExternalGrpcPort)
	if err := registry.Register(ctx, grpcInstanceID, cfg.Service.Name, env.ExternalGrpcPort, []string{TagGRPC}); err != nil {
		log.Fatalf("could not register grpc service: %v", err)
	}
	defer registry.Deregister(ctx, grpcInstanceID, cfg.Service.Name)

	go func() {
		for {
			if err := registry.ReportHealthyState(restInstanceID, cfg.Service.Name); err != nil {
				log.Printf("could not report rest healthy state: %v", err)
			}
			if err := registry.ReportHealthyState(grpcInstanceID, cfg.Service.Name); err != nil {
				log.Printf("could not report grpc healthy state: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	storages, err := mysql.New(env.DB)
	if err != nil {
		log.Fatal(err)
	}

	authGateway := authService.New(registry)
	services := service.New(storages.User, authGateway)

	go func() {
		server := rest.New(cfg.RestServer).Init(services.User, authGateway)
		if err := server.Run(); err != nil {
			log.Fatalf("Rest server err: %v", err)
		}
	}()

	go func() {
		server := grpc.New(cfg.GrpcServer).Init(services.User, authGateway)
		if err := server.Run(); err != nil {
			log.Fatalf("Grpc server err: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")
}