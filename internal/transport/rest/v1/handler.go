package v1

import (
	"github.com/labstack/echo/v4"
	"log"
)

type (
	UserService interface {
	}

	TokenService interface {
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
