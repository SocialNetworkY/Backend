package main

import (
	"github.com/SocialNetworkY/Backend/internal/post/repository"
	"gorm.io/driver/mysql"
	"log"

	"github.com/SocialNetworkY/Backend/internal/post/service"
	"github.com/SocialNetworkY/Backend/internal/post/transport/grpc"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	DB   string `env:"DB"`
	Port int    `env:"PORT"`
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

	services := service.New(repos.Post, repos.Tag, repos.Like, repos.Comment)

	if err := grpc.New(cfg.Port).Init(services.Post).Run(); err != nil {
		log.Fatalf("Grpc server err: %v", err)
	}
}
