package main

import (
	"fmt"
	"log"

	"github.com/SocialNetworkY/Backend/internal/post/elasticsearch"
	"github.com/SocialNetworkY/Backend/internal/post/gateway/report"
	"github.com/SocialNetworkY/Backend/internal/post/gateway/user"
	"github.com/SocialNetworkY/Backend/pkg/storage"
	"gorm.io/driver/mysql"

	"github.com/SocialNetworkY/Backend/internal/post/gateway/auth"
	"github.com/SocialNetworkY/Backend/internal/post/repository"
	"github.com/SocialNetworkY/Backend/internal/post/service"
	"github.com/SocialNetworkY/Backend/internal/post/transport/http"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	DB                       string   `env:"DB"`
	Port                     int      `env:"PORT"`
	BodyLimit                string   `env:"BODY_LIMIT"`
	AllowedOrigins           []string `env:"ALLOWED_ORIGINS" envSeparator:","`
	PostElasticSearchAddr    string   `env:"POST_ELASTICSEARCH_ADDR"`
	TagElasticSearchAddr     string   `env:"TAG_ELASTICSEARCH_ADDR"`
	CommentElasticSearchAddr string   `env:"COMMENT_ELASTICSEARCH_ADDR"`
	AuthServiceHttpAddr      string   `env:"AUTH_SERVICE_HTTP_ADDR"`
	AuthServiceGrpcAddr      string   `env:"AUTH_SERVICE_GRPC_ADDR"`
	UserServiceHttpAddr      string   `env:"USER_SERVICE_HTTP_ADDR"`
	UserServiceGrpcAddr      string   `env:"USER_SERVICE_GRPC_ADDR"`
	ReportServiceHttpAddr    string   `env:"REPORT_SERVICE_HTTP_ADDR"`
	ReportServiceGrpcAddr    string   `env:"REPORT_SERVICE_GRPC_ADDR"`
	StorageFolder            string   `env:"STORAGE_FOLDER" envDefault:"storage"`
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
	postSearch, err := elasticsearch.NewPost(cfg.PostElasticSearchAddr)
	if err != nil {
		log.Fatalf("Post Elasticsearch err: %v", err)
	}

	tagSearch, err := elasticsearch.NewTag(cfg.TagElasticSearchAddr)
	if err != nil {
		log.Fatalf("Tag Elasticsearch err: %v", err)
	}

	commentSearch, err := elasticsearch.NewComment(cfg.CommentElasticSearchAddr)
	if err != nil {
		log.Fatalf("Comment Elasticsearch err: %v", err)
	}

	repos, err := repository.New(mysql.Open(cfg.DB), postSearch, commentSearch, tagSearch)
	if err != nil {
		log.Fatal(err)
	}

	fileStorage, err := storage.NewLocalStorage(cfg.StorageFolder, fmt.Sprintf("http://localhost:%d/%s", cfg.Port, "storage"))
	if err != nil {
		log.Fatalf("Image storage err: %v", err)
	}

	authGateway := auth.New(cfg.AuthServiceHttpAddr, cfg.AuthServiceGrpcAddr)
	postGateway := user.New(cfg.UserServiceHttpAddr, cfg.UserServiceGrpcAddr)
	reportGateway := report.New(cfg.ReportServiceHttpAddr, cfg.ReportServiceGrpcAddr)
	services := service.New(repos.Post, repos.Tag, repos.Like, repos.Comment, reportGateway)

	if err := http.New(cfg.BodyLimit, cfg.AllowedOrigins, cfg.Port).Init(services.Post, services.Like, services.Comment, authGateway, postGateway, fileStorage).AddStaticFolder("storage", cfg.StorageFolder).Run(); err != nil {
		log.Fatalf("Http server err: %v", err)
	}
}
