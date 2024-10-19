package main

import (
	"fmt"
	"github.com/SocialNetworkY/Backend/internal/user/elasticsearch"
	"github.com/SocialNetworkY/Backend/internal/user/gateway/post"
	"github.com/SocialNetworkY/Backend/internal/user/gateway/report"
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
	DB                    string `env:"DB"`
	Port                  int    `env:"PORT"`
	UserElasticSearchAddr string `env:"USER_ELASTICSEARCH_ADDR"`
	BanElasticSearchAddr  string `env:"BAN_ELASTICSEARCH_ADDR"`
	AuthServiceHttpAddr   string `env:"AUTH_SERVICE_HTTP_ADDR"`
	AuthServiceGrpcAddr   string `env:"AUTH_SERVICE_GRPC_ADDR"`
	PostServiceHttpAddr   string `env:"POST_SERVICE_HTTP_ADDR"`
	PostServiceGrpcAddr   string `env:"POST_SERVICE_GRPC_ADDR"`
	ReportServiceHttpAddr string `env:"REPORT_SERVICE_HTTP_ADDR"`
	ReportServiceGrpcAddr string `env:"REPORT_SERVICE_GRPC_ADDR"`
	StorageFolder         string `env:"STORAGE_FOLDER" envDefault:"storage"`
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
	userElastic, err := elasticsearch.NewUser(cfg.UserElasticSearchAddr)
	if err != nil {
		log.Fatalf("Elasticsearch err: %v", err)
	}

	banElastic, err := elasticsearch.NewBan(cfg.BanElasticSearchAddr)
	if err != nil {
		log.Fatalf("Elasticsearch err: %v", err)
	}

	repos, err := repository.New(mysql.Open(cfg.DB), userElastic, banElastic)
	if err != nil {
		log.Fatal(err)
	}

	imageStorage, err := storage.NewLocalStorage(cfg.StorageFolder, fmt.Sprintf("http://localhost:%d/%s", cfg.Port, "storage"))
	if err != nil {
		log.Fatal(err)
	}

	authGateway := auth.New(cfg.AuthServiceHttpAddr, cfg.AuthServiceGrpcAddr)
	postGateway := post.New(cfg.PostServiceHttpAddr, cfg.PostServiceGrpcAddr)
	reportGateway := report.New(cfg.ReportServiceHttpAddr, cfg.ReportServiceGrpcAddr)
	services := service.New(repos.User, repos.Ban, imageStorage, authGateway, postGateway, reportGateway)

	if err := grpc.New(cfg.Port).Init(services.User, authGateway).Run(); err != nil {
		log.Fatalf("Grpc server err: %v", err)
	}
}
