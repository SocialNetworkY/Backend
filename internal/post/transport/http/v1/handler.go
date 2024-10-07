package v1

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
	"io"
	"log"
	"time"
)

type (
	PostService interface {
		Create(post *model.Post) error
		Find(id uint) (*model.Post, error)
		Update(post *model.Post) error
		Delete(id uint) error
	}

	AuthGateway interface {
		Authenticate(ctx context.Context, auth string) (uint, error)
	}

	UserGateway interface {
		GetUserRole(ctx context.Context, auth string, userID uint) (uint, error)
		IsUserBan(ctx context.Context, auth string, userID uint) (banned bool, reason string, expiredAt time.Time, err error)
	}

	FileStorage interface {
		UploadFile(file io.ReadSeeker, fileName string) (string, error)
	}

	Handler struct {
		ps PostService
		ag AuthGateway
		ug UserGateway
		fs FileStorage
	}
)

func New(ps PostService, ag AuthGateway, ug UserGateway, fs FileStorage) *Handler {
	return &Handler{
		ps: ps,
		ag: ag,
		ug: ug,
		fs: fs,
	}
}

func (h *Handler) Init(api *echo.Group) {
	log.Println("Initializing V1 api")
	v1 := api.Group("/v1")
	{
		h.initPostApi(v1)
	}
}
