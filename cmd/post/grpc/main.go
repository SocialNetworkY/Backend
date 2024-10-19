package main

import (
	"log"

	"github.com/SocialNetworkY/Backend/internal/post/elasticsearch"
	"github.com/SocialNetworkY/Backend/internal/post/gateway/report"
	"github.com/SocialNetworkY/Backend/internal/post/repository"
	"gorm.io/driver/mysql"

	"github.com/SocialNetworkY/Backend/internal/post/service"
	"github.com/SocialNetworkY/Backend/internal/post/transport/grpc"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	DB                       string `env:"DB"`
	Port                     int    `env:"PORT"`
	PostElasticSearchAddr    string `env:"POST_ELASTICSEARCH_ADDR"`
	TagElasticSearchAddr     string `env:"TAG_ELASTICSEARCH_ADDR"`
	CommentElasticSearchAddr string `env:"COMMENT_ELASTICSEARCH_ADDR"`
	ReportServiceHttpAddr    string `env:"REPORT_SERVICE_HTTP_ADDR"`
	ReportServiceGrpcAddr    string `env:"REPORT_SERVICE_GRPC_ADDR"`
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

	reportGateway := report.New(cfg.ReportServiceHttpAddr, cfg.ReportServiceGrpcAddr)
	services := service.New(repos.Post, repos.Tag, repos.Like, repos.Comment, reportGateway)

	if err := grpc.New(cfg.Port).Init(services.Post, services.Comment, services.Like).Run(); err != nil {
		log.Fatalf("Grpc server err: %v", err)
	}
}
