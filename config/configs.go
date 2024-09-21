package config

import (
	"github.com/lapkomo2018/goTwitterServices/internal/auth/transport/grpc"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/transport/rest"
	"github.com/lapkomo2018/goTwitterServices/pkg/jwt"
	"github.com/lapkomo2018/goTwitterServices/pkg/validation"
)

type AuthConfig struct {
	Service struct {
		Name string
	}
	RestServer rest.Config
	GrpcServer grpc.Config
	JWT        jwt.Config
	Validator  validation.Config
}
