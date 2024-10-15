package v1

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/auth/model"
	"github.com/SocialNetworkY/Backend/pkg/gen"
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
	user, err := h.as.Auth(r.GetAuthToken())
	if err != nil {
		return nil, err
	}

	return &gen.AuthenticateResponse{UserId: uint64(user.ID)}, nil
}

func (h *Handler) UpdateUsernameEmail(ctx context.Context, r *gen.UpdateUsernameEmailRequest) (*gen.UpdateUsernameEmailResponse, error) {
	userID := uint(r.GetUserId())

	if r.Username != "" {
		if err := h.us.ChangeUsername(userID, r.Username); err != nil {
			return &gen.UpdateUsernameEmailResponse{Success: false}, err
		}
	}

	if r.Email != "" {
		if err := h.us.ChangeEmail(userID, r.Email); err != nil {
			return &gen.UpdateUsernameEmailResponse{Success: false}, err
		}
	}

	return &gen.UpdateUsernameEmailResponse{Success: true}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, r *gen.DeleteUserRequest) (*gen.DeleteUserResponse, error) {
	userID := uint(r.GetUserId())
	if err := h.us.Delete(userID); err != nil {
		return &gen.DeleteUserResponse{Success: false}, err
	}

	return &gen.DeleteUserResponse{Success: true}, nil
}
