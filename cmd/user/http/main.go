package main

import (
	"fmt"
	"github.com/lapkomo2018/goTwitterServices/pkg/storage"
	"log"

	"github.com/lapkomo2018/goTwitterServices/internal/user/gateway/auth"
	"github.com/lapkomo2018/goTwitterServices/internal/user/repository/mysql"
	"github.com/lapkomo2018/goTwitterServices/internal/user/service"
	"github.com/lapkomo2018/goTwitterServices/internal/user/transport/http"

	"github.com/lapkomo2018/goTwitterServices/pkg/config"

	envCarlos "github.com/caarlos0/env/v6"
)

type Config struct {
	HttpServer http.Config
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

const (
	ImageFolder = "images"
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

	imageStorage, err := storage.NewLocalStorage(ImageFolder, fmt.Sprintf("http://localhost:%d/%s", env.Port, ImageFolder))
	if err != nil {
		log.Fatalf("Image storage err: %v", err)
	}

	authGateway := auth.New(env.AuthServiceHttpAddr, env.AuthServiceGrpcAddr)
	services := service.New(storages.User, storages.Ban, imageStorage, authGateway)

	if err := http.New(cfg.HttpServer, env.Port).Init(services.User, services.Ban, authGateway).AddStaticFolder(ImageFolder, ImageFolder).Run(); err != nil {
		log.Fatalf("Http server err: %v", err)
	}
}
