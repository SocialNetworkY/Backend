package main

import (
	"fmt"
	"github.com/SocialNetworkY/Backend/internal/user/gateway/post"
	"github.com/SocialNetworkY/Backend/internal/user/repository"
	"github.com/SocialNetworkY/Backend/pkg/storage"
	"gorm.io/driver/mysql"
	"log"

	"github.com/SocialNetworkY/Backend/internal/user/gateway/auth"
	"github.com/SocialNetworkY/Backend/internal/user/service"
	"github.com/SocialNetworkY/Backend/internal/user/transport/grpc"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	DB                  string `env:"DB"`
	Port                int    `env:"PORT"`
	AuthServiceHttpAddr string `env:"AUTH_SERVICE_HTTP_ADDR"`
	AuthServiceGrpcAddr string `env:"AUTH_SERVICE_GRPC_ADDR"`
	PostServiceHttpAddr string `env:"POST_SERVICE_HTTP_ADDR"`
	PostServiceGrpcAddr string `env:"POST_SERVICE_GRPC_ADDR"`
	StorageFolder       string `env:"STORAGE_FOLDER" envDefault:"storage"`
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

	imageStorage, err := storage.NewLocalStorage(cfg.StorageFolder, fmt.Sprintf("http://localhost:%d/%s", cfg.Port, "storage"))
	if err != nil {
		log.Fatal(err)
	}

	authGateway := auth.New(cfg.AuthServiceHttpAddr, cfg.AuthServiceGrpcAddr)
	postGateway := post.New(cfg.PostServiceHttpAddr, cfg.PostServiceGrpcAddr)
	services := service.New(repos.User, repos.Ban, imageStorage, authGateway, postGateway)

	if err := grpc.New(cfg.Port).Init(services.User, authGateway).Run(); err != nil {
		log.Fatalf("Grpc server err: %v", err)
	}
}
