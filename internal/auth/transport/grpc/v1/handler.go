package v1

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/model"
	"github.com/lapkomo2018/goTwitterServices/pkg/gen"
)

type (
	AuthenticationService interface {
		Auth(auth string) (*model.User, error)
	}

	UserService interface {
		ChangeEmail(id uint, email string) error
		ChangeUsername(id uint, username string) error
		Delete(id uint) error
	}

	Handler struct {
		gen.UnimplementedAuthServiceServer
		as AuthenticationService
		us UserService
	}
)

func New(as AuthenticationService, us UserService) *Handler {
	return &Handler{
		as: as,
		us: us,
	}
}

func (h *Handler) Authenticate(ctx context.Context, r *gen.AuthenticateRequest) (*gen.AuthenticateResponse, error) {
	user, err := h.getUserFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return &gen.AuthenticateResponse{UserId: uint64(user.ID)}, nil
}

// TODO: Fix this (Need update credentials for user from request id)
/*func (h *Handler) UpdateUsernameEmail(ctx context.Context, r *gen.UpdateUsernameEmailRequest) (*gen.UpdateUsernameEmailResponse, error) {
	user, err := h.getUserFromMetadata(ctx)
	if err != nil {
		return &gen.UpdateUsernameEmailResponse{Success: false}, err
	}

	if r.Username != "" {
		err = h.us.ChangeUsername(user.ID, r.Username)
		if err != nil {
			return &gen.UpdateUsernameEmailResponse{Success: false}, err
		}
	}

	if r.Email != "" {
		err = h.us.ChangeEmail(user.ID, r.Email)
		if err != nil {
			return &gen.UpdateUsernameEmailResponse{Success: false}, err
		}
	}

	return &gen.UpdateUsernameEmailResponse{Success: true}, nil
}*/

// TODO: Fix this (Need delete user from request id not from token)
/*func (h *Handler) DeleteUser(ctx context.Context, r *gen.DeleteUserRequest) (*gen.DeleteUserResponse, error) {
	user, err := h.getUserFromMetadata(ctx)
	if err != nil {
		return &gen.DeleteUserResponse{Success: false}, err
	}

	err = h.us.Delete(user.ID)
	if err != nil {
		return &gen.DeleteUserResponse{Success: false}, err
	}

	return &gen.DeleteUserResponse{Success: true}, nil
}*/
