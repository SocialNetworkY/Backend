package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig[T any]() (*T, error) {
	log.Println("Loading config...")

	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config T
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	log.Println("Loaded config")
	return &config, nil
}
