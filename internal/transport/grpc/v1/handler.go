package v1

import (
	"context"
	"github.com/lapkomo2018/goTwitterAuthService/internal/core"
	grpcService "github.com/lapkomo2018/goTwitterAuthService/pkg/grpc/auth"
)

type (
	AuthenticationService interface {
		Auth(auth string) (*core.User, error)
	}

	Handler struct {
		grpcService.UnimplementedAuthenticationServer
		authenticationService AuthenticationService
	}
)

func New(authenticationService AuthenticationService) *Handler {
	return &Handler{
		authenticationService: authenticationService,
	}
}

func (h *Handler) Authenticate(ctx context.Context, r *grpcService.AuthenticateRequest) (*grpcService.AuthenticateResponse, error) {
	user, err := h.authenticationService.Auth(r.GetAuth())
	if err != nil {
		return nil, err
	}
	return &grpcService.AuthenticateResponse{UserId: uint64(user.ID)}, nil
}
