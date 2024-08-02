package main

import (
	"github.com/lapkomo2018/goTwitterAuthService/internal/service"
	"github.com/lapkomo2018/goTwitterAuthService/internal/storage/mysql"
	grpcServer "github.com/lapkomo2018/goTwitterAuthService/internal/transport/grpc"
	restServer "github.com/lapkomo2018/goTwitterAuthService/internal/transport/rest"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const httpBodyLimit = 1024 * 1024 * 24

// @title           Twitter Auth Service
// @version         1.0
// @description     Bombaclac

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	storages, err := mysql.New("root:strongpass@tcp(localhost:3306)/main")
	if err != nil {
		log.Fatal(err)
	}

	services := service.New(storages.User, storages.RefreshToken)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		server := restServer.New(httpBodyLimit, []string{"*"}).Init(services.User, storages.RefreshToken)
		if err := server.Run(8080); err != nil {
			log.Fatalf("Rest server err: %v", err)
		}
	}()

	go func() {
		server := grpcServer.New(8081).Init()
		if err := server.Run(); err != nil {
			log.Fatalf("Grpc server err: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down servers...")
}
