package main

import (
	"log"
	"time"

	"github.com/SocialNetworkY/Backend/internal/auth/repository"
	"gorm.io/driver/mysql"

	"github.com/SocialNetworkY/Backend/internal/auth/gateway/user"
	"github.com/SocialNetworkY/Backend/internal/auth/service"
	"github.com/SocialNetworkY/Backend/internal/auth/transport/http"

	"github.com/SocialNetworkY/Backend/pkg/hash"
	"github.com/SocialNetworkY/Backend/pkg/jwt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	DB                  string        `env:"DB"`
	Port                int           `env:"PORT"`
	BodyLimit           string        `env:"BODY_LIMIT"`
	AllowedOrigins      []string      `env:"ALLOWED_ORIGINS" envSeparator:","`
	HashSalt            string        `env:"HASH_SALT"`
	JWTSecret           string        `env:"JWT_SECRET"`
	JWTDuration         time.Duration `env:"JWT_DURATION"`
	JWTRefreshSecret    string        `env:"JWT_REFRESH_SECRET"`
	JWTRefreshDuration  time.Duration `env:"JWT_REFRESH_DURATION"`
	UserServiceHttpAddr string        `env:"USER_SERVICE_HTTP_ADDR"`
	UserServiceGrpcAddr string        `env:"USER_SERVICE_GRPC_ADDR"`
}

var (
	cfg = Config{}
)

func init() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
}

func main() {
	repos, err := repository.New(mysql.Open(cfg.DB))
	if err != nil {
		log.Fatal(err)
	}

	hasher := hash.NewSHA1Hasher(cfg.HashSalt)
	tokenManager := jwt.NewManager(cfg.JWTDuration, cfg.JWTRefreshDuration, cfg.JWTSecret, cfg.JWTRefreshSecret)
	userGateway := user.New(cfg.UserServiceHttpAddr, cfg.UserServiceGrpcAddr)
	services := service.New(repos.User, repos.RefreshToken, repos.ActivationToken, tokenManager, hasher, userGateway)

	if err := http.New(cfg.BodyLimit, cfg.AllowedOrigins, cfg.Port).Init(services.User, services.Tokens, services.Authentication, cfg.JWTRefreshDuration).Run(); err != nil {
		log.Fatalf("Rest server err: %v", err)
	}
}
