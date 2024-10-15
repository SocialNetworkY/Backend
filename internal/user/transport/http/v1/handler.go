package v1

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/user/model"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"time"
)

type (
	UserService interface {
		Find(id uint) (*model.User, error)
		FindByUsername(username string) (*model.User, error)
		FindByEmail(email string) (*model.User, error)
		FindByNickname(nickname string, skip, limit int) ([]*model.User, error)
		ChangeEmail(id uint, email string) error
		ChangeUsername(id uint, username string) error
		ChangeNickname(id uint, nickname string) error
		ChangeAvatar(id uint, file io.ReadSeeker) error
		Delete(id uint) error
	}

	BanService interface {
		Ban(userID, adminID uint, reason string, duration time.Duration) error
		Unban(banID, adminID uint, reason string) error
		Find(id uint) (*model.Ban, error)
		FindByUser(userID uint, skip, limit int) ([]*model.Ban, error)
	}

	AuthGateway interface {
		Authenticate(ctx context.Context, auth string) (uint, error)
	}

	Handler struct {
		us UserService
		bs BanService
		ag AuthGateway
	}
)

func New(us UserService, bs BanService, ag AuthGateway) *Handler {
	return &Handler{
		us: us,
		ag: ag,
		bs: bs,
	}
}

func (h *Handler) Init(api *echo.Group) {
	log.Println("Initializing V1 api")
	v1 := api.Group("/v1")
	{
		h.initUserApi(v1)
		h.initAdminApi(v1)
	}
}
