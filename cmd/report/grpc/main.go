package main

import (
	"log"

	"github.com/SocialNetworkY/Backend/internal/report/elasticsearch"
	"github.com/SocialNetworkY/Backend/internal/report/gateway/post"

	"github.com/SocialNetworkY/Backend/internal/report/repository"
	"github.com/SocialNetworkY/Backend/internal/report/service"
	"github.com/SocialNetworkY/Backend/internal/report/transport/grpc"

	"github.com/caarlos0/env/v6"
	"gorm.io/driver/mysql"
)

type Config struct {
	DB                      string `env:"DB"`
	Port                    int    `env:"PORT"`
	ReportElasticSearchAddr string `env:"REPORT_ELASTICSEARCH_ADDR"`
	PostServiceHttpAddr     string `env:"POST_SERVICE_HTTP_ADDR"`
	PostServiceGrpcAddr     string `env:"POST_SERVICE_GRPC_ADDR"`
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
	reportSearch, err := elasticsearch.NewReport(cfg.ReportElasticSearchAddr)
	if err != nil {
		log.Fatal(err)
	}

	repos, err := repository.New(mysql.Open(cfg.DB), reportSearch)
	if err != nil {
		log.Fatal(err)
	}

	postGateway := post.New(cfg.PostServiceHttpAddr, cfg.PostServiceGrpcAddr)
	services := service.New(repos.Report, postGateway)

	if err := grpc.New(cfg.Port).Init(services.Report).Run(); err != nil {
		log.Fatalf("Grpc server err: %v", err)
	}
}
