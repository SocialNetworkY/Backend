package main

import (
	"github.com/SocialNetworkY/Backend/internal/post/gateway/user"
	"github.com/SocialNetworkY/Backend/internal/post/repository"
	"gorm.io/driver/mysql"
	"log"

	"github.com/SocialNetworkY/Backend/internal/post/gateway/auth"
	"github.com/SocialNetworkY/Backend/internal/post/service"
	"github.com/SocialNetworkY/Backend/internal/post/transport/grpc"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	DB                  string `env:"DB"`
	Port                int    `env:"PORT"`
	AuthServiceHttpAddr string `env:"AUTH_SERVICE_HTTP_ADDR"`
	AuthServiceGrpcAddr string `env:"AUTH_SERVICE_GRPC_ADDR"`
	UserServiceHttpAddr string `env:"USER_SERVICE_HTTP_ADDR"`
	UserServiceGrpcAddr string `env:"USER_SERVICE_GRPC_ADDR"`
}

var (
	cfg = &Config{}
)

func init() {
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
}

func main() {
	repos, err := repository.New(mysql.Open(cfg.DB))
	if err != nil {
		log.Fatal(err)
	}

	_ = auth.New(cfg.AuthServiceHttpAddr, cfg.AuthServiceGrpcAddr)
	_ = user.New(cfg.UserServiceHttpAddr, cfg.UserServiceGrpcAddr)
	_ = service.New(repos.Post, repos.Tag, repos.Like, repos.Comment)

	if err := grpc.New(cfg.Port).Init().Run(); err != nil {
		log.Fatalf("Grpc server err: %v", err)
	}
}
