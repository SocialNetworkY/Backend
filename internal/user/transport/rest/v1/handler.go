package v1

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"log"
)

type (
	UserService interface {
		Find(id uint) (*model.User, error)
		FindByUsername(username string) (*model.User, error)
		FindByEmail(email string) (*model.User, error)
		FindByNickname(nickname string) ([]*model.User, error)
		ChangeEmail(id uint, auth, email string) error
		ChangeUsername(id uint, auth, username string) error
		ChangeNickname(id uint, nickname string) error
		ChangeAvatar(id uint, avatar string) error
		Delete(id uint, auth string) error
	}

	AuthGateway interface {
		CheckAuth(ctx context.Context, auth string) (uint, error)
	}

	Handler struct {
		us UserService
		ag AuthGateway
	}
)

const (
	userLocals = "user"
)

func New(us UserService, ag AuthGateway) *Handler {
	return &Handler{
		us: us,
		ag: ag,
	}
}

func (h *Handler) Init(api *echo.Group) {
	log.Println("Initializing V1 api")
	v1 := api.Group("/v1")
	{
		h.initUserApi(v1)
	}
}
