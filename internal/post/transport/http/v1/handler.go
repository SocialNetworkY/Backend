package v1

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/post/model"
	"github.com/labstack/echo/v4"
	"io"
	"log"
)

type (
	PostService interface {
		Create(post *model.Post) error
		Update(post *model.Post) error
		Delete(id uint) error
		Find(id uint) (*model.Post, error)
		FindSome(skip, limit int) ([]*model.Post, error)
		FindByUser(userID uint, skip, limit int) ([]*model.Post, error)
	}

	LikeService interface {
		LikePost(postID, userID uint) error
		UnlikePost(postID, userID uint) error
	}

	CommentService interface {
		Find(id uint) (*model.Comment, error)
		FindByPost(postID uint, skip, limit int) ([]*model.Comment, error)
		CommentPost(postID, userID uint, content string) error
		Edit(id, userID uint, content string) error
		Delete(id uint) error
	}

	AuthGateway interface {
		Authenticate(ctx context.Context, auth string) (uint, error)
	}

	UserGateway interface {
		UserInfo(ctx context.Context, userID uint) (*model.User, error)
	}

	FileStorage interface {
		UploadFile(file io.ReadSeeker, fileName string) (string, error)
	}

	Handler struct {
		ps PostService
		ls LikeService
		cs CommentService
		ag AuthGateway
		ug UserGateway
		fs FileStorage
	}
)

func New(ps PostService, ls LikeService, cs CommentService, ag AuthGateway, ug UserGateway, fs FileStorage) *Handler {
	return &Handler{
		ps: ps,
		ls: ls,
		cs: cs,
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
		h.initCommentsApi(v1)
	}
}
