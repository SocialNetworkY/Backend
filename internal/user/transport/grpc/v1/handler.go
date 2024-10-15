package v1

import (
	"context"
	"errors"
	"github.com/SocialNetworkY/Backend/internal/user/model"
	"github.com/SocialNetworkY/Backend/pkg/gen"
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

func (h *Handler) CreateUser(ctx context.Context, r *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {
	user, err := h.us.Create(uint(r.GetUserId()), uint(r.GetRole()), r.GetUsername(), r.GetEmail())
	if err != nil {
		return nil, err
	}

	return &gen.CreateUserResponse{
		UserId: uint64(user.ID),
	}, nil
}

func (h *Handler) UserInfo(ctx context.Context, r *gen.UserInfoRequest) (*gen.UserInfoResponse, error) {
	userID := uint(r.GetUserId())
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}

	user, err := h.us.Find(userID)
	if err != nil {
		return nil, err
	}

	banReason := ""
	banExpiredAt := ""
	if user.Banned {
		banReason = user.ActiveBan.BanReason
		banExpiredAt = user.ActiveBan.ExpiredAt.Format(time.RFC3339)
	}

	return &gen.UserInfoResponse{
		UserId:       uint64(user.ID),
		Role:         uint64(user.Role),
		Banned:       user.Banned,
		BanReason:    banReason,
		BanExpiredAt: banExpiredAt,
	}, nil
}
