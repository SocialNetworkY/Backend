package v1

import (
	"context"
	"github.com/lapkomo2018/goTwitterAuthService/pkg/gen"
	"github.com/lapkomo2018/goTwitterAuthService/pkg/model"
)

type (
	AuthenticationService interface {
		Auth(auth string) (*model.User, error)
	}

	Handler struct {
		gen.UnimplementedAuthenticationServer
		authenticationService AuthenticationService
	}
)

func New(authenticationService AuthenticationService) *Handler {
	return &Handler{
		authenticationService: authenticationService,
	}
}

func (h *Handler) Authenticate(ctx context.Context, r *gen.AuthenticateRequest) (*gen.AuthenticateResponse, error) {
	user, err := h.authenticationService.Auth(r.GetAuth())
	if err != nil {
		return nil, err
	}
	return &gen.AuthenticateResponse{UserId: uint64(user.ID)}, nil
}
