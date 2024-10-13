package main

import (
	"fmt"
	"github.com/SocialNetworkY/Backend/internal/user/repository"
	"github.com/SocialNetworkY/Backend/pkg/storage"
	"gorm.io/driver/mysql"
	"log"

	"github.com/SocialNetworkY/Backend/internal/user/gateway/auth"
	"github.com/SocialNetworkY/Backend/internal/user/service"
	"github.com/SocialNetworkY/Backend/internal/user/transport/http"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	DB                  string   `env:"DB"`
	Port                int      `env:"PORT"`
	BodyLimit           string   `env:"BODY_LIMIT"`
	AllowedOrigins      []string `env:"ALlOWED_ORIGINS" envSeparator:","`
	AuthServiceHttpAddr string   `env:"AUTH_SERVICE_HTTP_ADDR"`
	AuthServiceGrpcAddr string   `env:"AUTH_SERVICE_GRPC_ADDR"`
	StorageFolder       string   `env:"STORAGE_FOLDER" envDefault:"storage"`
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
		log.Fatalf("Image storage err: %v", err)
	}

	authGateway := auth.New(cfg.AuthServiceHttpAddr, cfg.AuthServiceGrpcAddr)
	services := service.New(repos.User, repos.Ban, imageStorage, authGateway)

	if err := http.New(cfg.BodyLimit, cfg.AllowedOrigins, cfg.Port).Init(services.User, services.Ban, authGateway).AddStaticFolder("storage", cfg.StorageFolder).Run(); err != nil {
		log.Fatalf("Http server err: %v", err)
	}
}
