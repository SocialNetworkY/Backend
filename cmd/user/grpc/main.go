package main

import (
	"log"

	"github.com/lapkomo2018/goTwitterServices/internal/user/gateway/auth"
	"github.com/lapkomo2018/goTwitterServices/internal/user/repository/mysql"
	"github.com/lapkomo2018/goTwitterServices/internal/user/service"
	"github.com/lapkomo2018/goTwitterServices/internal/user/transport/grpc"

	"github.com/lapkomo2018/goTwitterServices/pkg/config"

	envCarlos "github.com/caarlos0/env/v6"
)

type Config struct {
	GrpcServer grpc.Config
}

type Env struct {
	DB                  string `env:"DB"`
	Port                int    `env:"PORT"`
	AuthServiceHttpAddr string `env:"AUTH_SERVICE_HTTP_ADDR"`
	AuthServiceGrpcAddr string `env:"AUTH_SERVICE_GRPC_ADDR"`
}

var (
	cfg *Config
	env = &Env{}
)

func init() {
	var err error
	if err := envCarlos.Parse(env); err != nil {
		log.Fatal(err)
	}

	cfg, err = config.LoadConfig[Config]("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	storages, err := mysql.New(env.DB)
	if err != nil {
		log.Fatal(err)
	}

	authGateway := auth.New(env.AuthServiceHttpAddr, env.AuthServiceGrpcAddr)
	services := service.New(storages.User, authGateway)

	if err := grpc.New(cfg.GrpcServer, env.Port).Init(services.User, authGateway).Run(); err != nil {
		log.Fatalf("Grpc server err: %v", err)
	}
}