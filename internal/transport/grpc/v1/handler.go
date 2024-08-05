package v1

import (
	"context"
	"github.com/lapkomo2018/goTwitterAuthService/internal/core"
	grpcAuthService "github.com/lapkomo2018/goTwitterAuthService/pkg/grpc/auth"
)

type (
	AuthenticationService interface {
		Auth(auth string) (*core.User, error)
	}

	Handler struct {
		grpcAuthService.UnimplementedAuthenticationServer
		authenticationService AuthenticationService
	}
)

func New(authenticationService AuthenticationService) *Handler {
	return &Handler{
		authenticationService: authenticationService,
	}
}

func (h *Handler) Authenticate(ctx context.Context, r *grpcAuthService.AuthenticateRequest) (*grpcAuthService.AuthenticateResponse, error) {
	user, err := h.authenticationService.Auth(r.GetAuth())
	if err != nil {
		return nil, err
	}
	return &grpcAuthService.AuthenticateResponse{UserId: uint64(user.ID)}, nil
}
