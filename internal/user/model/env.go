package model

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Env struct {
	DB               string
	DiscoveryAddr    string
	ExternalRestPort int
	ExternalGrpcPort int
}

const (
	envDB               = "DB"
	envDiscoveryAddr    = "DISCOVERY_ADDR"
	envExternalRestPort = "EXTERNAL_REST_PORT"
	envExternalGrpcPort = "EXTERNAL_GRPC_PORT"
)

func ParseEnv() (*Env, error) {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file, proceeding without it: %v", err)
	}

	restPort, err := strconv.Atoi(os.Getenv(envExternalRestPort))
	if err != nil {
		return nil, err
	}

	grpcPort, err := strconv.Atoi(os.Getenv(envExternalGrpcPort))
	if err != nil {
		return nil, err
	}

	return &Env{
		DB:               os.Getenv(envDB),
		DiscoveryAddr:    os.Getenv(envDiscoveryAddr),
		ExternalRestPort: restPort,
		ExternalGrpcPort: grpcPort,
	}, nil
}
