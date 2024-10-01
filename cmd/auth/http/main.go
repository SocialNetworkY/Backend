package main

import (
	"log"

	"github.com/lapkomo2018/goTwitterServices/internal/auth/gateway/user"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/repository/mysql"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/service"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/transport/rest"

	"github.com/lapkomo2018/goTwitterServices/pkg/config"
	"github.com/lapkomo2018/goTwitterServices/pkg/hash"
	"github.com/lapkomo2018/goTwitterServices/pkg/jwt"
	"github.com/lapkomo2018/goTwitterServices/pkg/validation"

	envCarlos "github.com/caarlos0/env/v6"
)

type Config struct {
	HttpServer rest.Config
	JWT        jwt.Config
	Validator  validation.Config
}

type Env struct {
	DB                  string `env:"DB"`
	Port                int    `env:"PORT"`
	HashSalt            string `env:"HASH_SALT"`
	JWTSecret           string `env:"JWT_SECRET"`
	JWTRefreshSecret    string `env:"JWT"`
	UserServiceHttpAddr string `env:"USER_SERVICE_HTTP_ADDR"`
	UserServiceGrpcAddr string `env:"USER_SERVICE_GRPC_ADDR"`
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

	validator, err := validation.NewValidator(cfg.Validator)
	if err != nil {
		log.Fatal(err)
	}

	hasher := hash.NewSHA1Hasher(env.HashSalt)
	tokenManager := jwt.NewManager(cfg.JWT, env.JWTSecret, env.JWTRefreshSecret)
	userGateway := user.New(env.UserServiceHttpAddr, env.UserServiceGrpcAddr)
	services := service.New(storages.User, storages.RefreshToken, storages.ActivationToken, tokenManager, hasher, userGateway)

	if err := rest.New(cfg.HttpServer, env.Port).Init(services.User, services.Tokens, services.Authentication, validator, cfg.JWT.RefreshDuration).Run(); err != nil {
		log.Fatalf("Rest server err: %v", err)
	}
}
