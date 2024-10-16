package main

import (
	"github.com/SocialNetworkY/Backend/internal/report/gateway/post"
	"github.com/SocialNetworkY/Backend/internal/report/transport/http"
	"log"

	"github.com/SocialNetworkY/Backend/internal/report/gateway/auth"
	"github.com/SocialNetworkY/Backend/internal/report/gateway/user"
	"github.com/SocialNetworkY/Backend/internal/report/repository"
	"github.com/SocialNetworkY/Backend/internal/report/service"
	"github.com/caarlos0/env/v6"
	"gorm.io/driver/mysql"
)

type Config struct {
	DB                  string   `env:"DB"`
	Port                int      `env:"PORT"`
	BodyLimit           string   `env:"BODY_LIMIT"`
	AllowedOrigins      []string `env:"ALlOWED_ORIGINS" envSeparator:","`
	AuthServiceHttpAddr string   `env:"AUTH_SERVICE_HTTP_ADDR"`
	AuthServiceGrpcAddr string   `env:"AUTH_SERVICE_GRPC_ADDR"`
	UserServiceHttpAddr string   `env:"USER_SERVICE_HTTP_ADDR"`
	UserServiceGrpcAddr string   `env:"USER_SERVICE_GRPC_ADDR"`
	PostServiceHttpAddr string   `env:"POST_SERVICE_HTTP_ADDR"`
	PostServiceGrpcAddr string   `env:"POST_SERVICE_GRPC_ADDR"`
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

	authGateway := auth.New(cfg.AuthServiceHttpAddr, cfg.AuthServiceGrpcAddr)
	userGateway := user.New(cfg.UserServiceHttpAddr, cfg.UserServiceGrpcAddr)
	postGateway := post.New(cfg.PostServiceHttpAddr, cfg.PostServiceGrpcAddr)
	services := service.New(repos.Report, postGateway)

	if err := http.New(cfg.BodyLimit, cfg.AllowedOrigins, cfg.Port).Init(services.Report, authGateway, userGateway).Run(); err != nil {
		log.Fatalf("Http server err: %v", err)
	}
}
