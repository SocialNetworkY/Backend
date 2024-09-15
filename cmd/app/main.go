package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lapkomo2018/goTwitterAuthService/config"
	"github.com/lapkomo2018/goTwitterAuthService/internal/service"
	"github.com/lapkomo2018/goTwitterAuthService/internal/storage/mysql"
	"github.com/lapkomo2018/goTwitterAuthService/internal/transport/grpc"
	"github.com/lapkomo2018/goTwitterAuthService/internal/transport/rest"
	"github.com/lapkomo2018/goTwitterAuthService/pkg/discovery"
	"github.com/lapkomo2018/goTwitterAuthService/pkg/discovery/consul"
	"github.com/lapkomo2018/goTwitterAuthService/pkg/hash"
	"github.com/lapkomo2018/goTwitterAuthService/pkg/jwt"
	"github.com/lapkomo2018/goTwitterAuthService/pkg/validation"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cfg *config.Config

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file, proceeding without it: %v", err)
	}

	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
}

// @title           Twitter Auth Service
// @version         1.0
// @description     Bombaclac

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	registry, err := consul.NewRegistry(cfg.Discovery.Address)
	if err != nil {
		log.Fatalf("could not create registry: %v", err)
	}
	ctx := context.Background()
	hostPort := fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.RestServer.Port)
	instanceID := discovery.GenerateInstanceID(cfg.Service.Name, hostPort)
	if err := registry.Register(ctx, instanceID, cfg.Service.Name, hostPort); err != nil {
		log.Fatalf("could not register service: %v", err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, cfg.Service.Name); err != nil {
				log.Printf("could not report healthy state: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, cfg.Service.Name)

	storages, err := mysql.New(os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
	}
	validator, err := validation.NewValidator(cfg.Validator)
	if err != nil {
		log.Fatal(err)
	}
	hasher := hash.NewSHA1Hasher(cfg.Hash)
	tokenManager := jwt.NewManager(cfg.JWT)
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
