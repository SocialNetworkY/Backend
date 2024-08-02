package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	RestServer struct {
		Port           int
		BodyLimit      int
		AllowedOrigins []string
	}

	GrpcServer struct {
		Port int
	}

	JWT struct {
		SecretKey        string
		TokenDuration    time.Duration
		RefreshSecretKey string
		RefreshDuration  time.Duration
	}
}

func LoadConfig() (*Config, error) {
	log.Println("Loading config...")

	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	log.Println("Loaded config")
	return &config, nil
}
