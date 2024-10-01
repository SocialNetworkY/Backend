package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig[T any](path string) (*T, error) {
	log.Println("Loading config...")

	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config T
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
