package v1

import (
	"github.com/labstack/echo/v4"
	"log"
)

type (
	UserService interface {
		Register(username, email, password string) (string, string, error)
		Login(login, password string) (string, string, error)
	}

	TokenService interface {
		Generate(userID uint) (string, string, error)
		Verify(accessToken string) (userID uint, err error)
		VerifyRefreshToken(refreshToken string) (userID uint, err error)
	}

	Handler struct {
		userService  UserService
		tokenService TokenService
	}
)

func New(userService UserService, tokenService TokenService) *Handler {
	return &Handler{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (h *Handler) Init(api *echo.Group) {
	log.Println("Initializing V1 api")
	v1 := api.Group("/v1")
	{
		h.initUserApi(v1)
		h.initTokenApi(v1)
	}
}
