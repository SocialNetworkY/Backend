package v1

import (
	"context"
	"errors"
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"github.com/lapkomo2018/goTwitterServices/pkg/constant"
	"github.com/lapkomo2018/goTwitterServices/pkg/gen"
	"time"
)

type (
	UserService interface {
		Create(id, role uint, username, email string) (*model.User, error)
		Find(id uint) (*model.User, error)
	}

	AuthGateway interface {
		Authenticate(ctx context.Context, auth string) (uint, error)
	}

	Handler struct {
		gen.UnimplementedUserServiceServer
		us UserService
		ag AuthGateway
	}
)

func New(us UserService, ag AuthGateway) *Handler {
	return &Handler{
		us: us,
		ag: ag,
	}
}

func (h *Handler) GetUserRole(ctx context.Context, r *gen.GetUserRoleRequest) (*gen.GetUserRoleResponse, error) {
	requester, err := h.getRequesterFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	if !(requester.Role >= constant.RoleAdminLvl1 || requester.ID == uint(r.GetUserId())) {
		return nil, errors.New("you are not allowed to get user role")
	}

	user, err := h.us.Find(uint(r.GetUserId()))
	if err != nil {
		return nil, err
	}

	return &gen.GetUserRoleResponse{
		Role: uint64(user.Role),
	}, nil
}

func (h *Handler) IsUserBan(ctx context.Context, r *gen.IsUserBanRequest) (*gen.IsUserBanResponse, error) {
	requester, err := h.getRequesterFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	if requester.ID != uint(r.GetUserId()) && !requester.Admin {
		return nil, errors.New("you are not allowed to check user ban")
	}

	user, err := h.us.Find(uint(r.GetUserId()))
	if err != nil {
		return nil, err
	}

	if !user.Banned {
		return &gen.IsUserBanResponse{
			Banned: user.Banned,
		}, nil
	}

	return &gen.IsUserBanResponse{
		Banned:    user.Banned,
		Reason:    user.ActiveBan.BanReason,
		ExpiredAt: user.ActiveBan.ExpiredAt.Format(time.RFC3339),
	}, nil
}

func (h *Handler) CreateUser(ctx context.Context, r *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {
	requesterID, err := h.getRequesterIDFromMetadata(ctx)

	if !(requesterID == uint(r.GetUserId()) && r.GetRole() < constant.RoleAdminLvl1) {
		return nil, errors.New("you are not allowed to create user")
	}

	user, err := h.us.Create(uint(r.GetUserId()), uint(r.GetRole()), r.GetUsername(), r.GetEmail())
	if err != nil {
		return nil, err
	}

	return &gen.CreateUserResponse{
		UserId: uint64(user.ID),
	}, nil
}
