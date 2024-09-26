package v1

import (
	"context"
	"errors"
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"github.com/lapkomo2018/goTwitterServices/pkg/gen"
)

type (
	UserService interface {
		Create(id, role uint, username, email string) (*model.User, error)
		Find(id uint) (*model.User, error)
	}

	AuthGateway interface {
		CheckAuth(ctx context.Context, auth string) (uint, error)
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
	requester, err := h.getUserFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	if !(requester.Role >= model.RoleAdminLvl1 || requester.ID == uint(r.GetUserId())) {
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

func (h *Handler) CreateUser(ctx context.Context, r *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {
	requesterID, err := h.getUserIDFromMetadata(ctx)

	if !(requesterID == uint(r.GetUserId()) && r.GetRole() < model.RoleAdminLvl1) {
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
