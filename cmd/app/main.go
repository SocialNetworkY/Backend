package main

import (
	"github.com/joho/godotenv"
	"github.com/lapkomo2018/goTwitterAuthService/config"
	"github.com/lapkomo2018/goTwitterAuthService/internal/service"
	"github.com/lapkomo2018/goTwitterAuthService/internal/storage/mysql"
	grpcServer "github.com/lapkomo2018/goTwitterAuthService/internal/transport/grpc"
	restServer "github.com/lapkomo2018/goTwitterAuthService/internal/transport/rest"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var cfg *config.Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
}

// @title           Twitter Auth Service
// @version         1.0
// @description     Bombaclac

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	storages, err := mysql.New(os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
	}

	services := service.New(storages.User, storages.RefreshToken)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		server := restServer.New(cfg.RestServer.BodyLimit, cfg.RestServer.AllowedOrigins).Init(services.User, storages.RefreshToken)
		if err := server.Run(cfg.RestServer.Port); err != nil {
			log.Fatalf("Rest server err: %v", err)
		}
	}()

	go func() {
		server := grpcServer.New(cfg.GrpcServer.Port).Init()
		if err := server.Run(); err != nil {
			log.Fatalf("Grpc server err: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down servers...")
}
