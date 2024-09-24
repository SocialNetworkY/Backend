package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lapkomo2018/goTwitterServices/config"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/model"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/repository/mysql"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/service"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/transport/grpc"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/transport/rest"
	"github.com/lapkomo2018/goTwitterServices/pkg/discovery"
	"github.com/lapkomo2018/goTwitterServices/pkg/discovery/consul"
	"github.com/lapkomo2018/goTwitterServices/pkg/hash"
	"github.com/lapkomo2018/goTwitterServices/pkg/jwt"
	"github.com/lapkomo2018/goTwitterServices/pkg/validation"
)

const (
	TagRest = "rest"
	TagGRPC = "grpc"
)

var (
	cfg *config.AuthConfig
	env *model.Env
)

func init() {
	var err error

	env, err = model.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err = config.LoadConfig[config.AuthConfig]()
	if err != nil {
		log.Fatal(err)
	}
}

// @title           Twitter Auth Service
// @version         1.0
// @description     Twitter Auth Service

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host      localhost:8080
// @BasePath  /api/v1
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
	validator, err := validation.NewValidator(cfg.Validator)
	if err != nil {
		log.Fatal(err)
	}
	hasher := hash.NewSHA1Hasher(env.HashSalt)
	tokenManager := jwt.NewManager(cfg.JWT, env.JWTSecret, env.JWTRefreshSecret)
	services := service.New(storages.User, storages.RefreshToken, storages.ActivationToken, tokenManager, hasher)

	go func() {
		server := rest.New(cfg.RestServer).Init(services.User, services.Tokens, services.Authentication, validator, cfg.JWT.RefreshDuration)
		if err := server.Run(); err != nil {
			log.Fatalf("Rest server err: %v", err)
		}
	}()

	go func() {
		server := grpc.New(cfg.GrpcServer).Init(services.Authentication)
		if err := server.Run(); err != nil {
			log.Fatalf("Grpc server err: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")
}
