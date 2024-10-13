package v1

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/auth/model"
	"github.com/SocialNetworkY/Backend/pkg/constant"
	"github.com/SocialNetworkY/Backend/pkg/gen"
)

type (
	AuthenticationService interface {
		Auth(auth string) (*model.User, error)
	}

	UserGateway interface {
		GetUserRole(ctx context.Context, auth string, userID uint) (uint, error)
	}

	UserService interface {
		ChangeEmail(id uint, email string) error
		ChangeUsername(id uint, username string) error
		Delete(id uint) error
	}

	Handler struct {
		gen.UnimplementedAuthServiceServer
		as AuthenticationService
		ug UserGateway
		us UserService
	}
)

func New(as AuthenticationService, us UserService, ug UserGateway) *Handler {
	return &Handler{
		as: as,
		ug: ug,
		us: us,
	}
}

func (h *Handler) Authenticate(ctx context.Context, r *gen.AuthenticateRequest) (*gen.AuthenticateResponse, error) {
	user, _, err := h.getRequesterFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return &gen.AuthenticateResponse{UserId: uint64(user.ID)}, nil
}

func (h *Handler) UpdateUsernameEmail(ctx context.Context, r *gen.UpdateUsernameEmailRequest) (*gen.UpdateUsernameEmailResponse, error) {
	requester, auth, err := h.getRequesterFromMetadata(ctx)
	if err != nil {
		return &gen.UpdateUsernameEmailResponse{Success: false}, err
	}

	requesterRole, err := h.ug.GetUserRole(ctx, auth, requester.ID)
	if err != nil {
		return &gen.UpdateUsernameEmailResponse{Success: false}, err
	}

	userID := uint(r.GetUserId())

	if requester.ID != userID && requesterRole < constant.RoleAdminLvl1 {
		return &gen.UpdateUsernameEmailResponse{Success: false}, nil
	}

	if r.Username != "" {
		err = h.us.ChangeUsername(userID, r.Username)
		if err != nil {
			return &gen.UpdateUsernameEmailResponse{Success: false}, err
		}
	}

	if r.Email != "" {
		err = h.us.ChangeEmail(userID, r.Email)
		if err != nil {
			return &gen.UpdateUsernameEmailResponse{Success: false}, err
		}
	}

	return &gen.UpdateUsernameEmailResponse{Success: true}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, r *gen.DeleteUserRequest) (*gen.DeleteUserResponse, error) {
	requester, auth, err := h.getRequesterFromMetadata(ctx)
	if err != nil {
		return &gen.DeleteUserResponse{Success: false}, err
	}

	requesterRole, err := h.ug.GetUserRole(ctx, auth, requester.ID)
	if err != nil {
		return &gen.DeleteUserResponse{Success: false}, err
	}

	userID := uint(r.GetUserId())

	if requester.ID != userID && requesterRole < constant.RoleAdminLvl3 {
		return &gen.DeleteUserResponse{Success: false}, nil
	}

	err = h.us.Delete(userID)
	if err != nil {
		return &gen.DeleteUserResponse{Success: false}, err
	}

	return &gen.DeleteUserResponse{Success: true}, nil
}
