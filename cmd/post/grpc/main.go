package main

import (
	"github.com/SocialNetworkY/Backend/internal/post/gateway/report"
	"github.com/SocialNetworkY/Backend/internal/post/repository"
	"gorm.io/driver/mysql"
	"log"

	"github.com/SocialNetworkY/Backend/internal/post/service"
	"github.com/SocialNetworkY/Backend/internal/post/transport/grpc"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	DB                    string `env:"DB"`
	Port                  int    `env:"PORT"`
	ReportServiceHttpAddr string `env:"REPORT_SERVICE_HTTP_ADDR"`
	ReportServiceGrpcAddr string `env:"REPORT_SERVICE_GRPC_ADDR"`
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

	reportGateway := report.New(cfg.ReportServiceHttpAddr, cfg.ReportServiceGrpcAddr)
	services := service.New(repos.Post, repos.Tag, repos.Like, repos.Comment, reportGateway)

	if err := grpc.New(cfg.Port).Init(services.Post, services.Comment, services.Like).Run(); err != nil {
		log.Fatalf("Grpc server err: %v", err)
	}
}
