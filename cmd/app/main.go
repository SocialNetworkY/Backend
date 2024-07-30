package main

import (
	grpcServer "github.com/lapkomo2018/goPetProjectServiceCollector/internal/transport/grpc/server"
	restServer "github.com/lapkomo2018/goPetProjectServiceCollector/internal/transport/rest/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const httpBodyLimit = 1024 * 1024 * 24

func startRest() {
	server := restServer.New(8080, httpBodyLimit, []string{"*"}).Init()
	if err := server.Run(); err != nil {
		log.Fatalf("Rest server err: %v", err)
	}
}

func startGrpc() {
	server := grpcServer.New(8081).Init()
	if err := server.Run(); err != nil {
		log.Fatalf("Grpc server err: %v", err)
	}
}

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go startRest()
	go startGrpc()

	<-quit
	log.Println("Shutting down servers...")
}
