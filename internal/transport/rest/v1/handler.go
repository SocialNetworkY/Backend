package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterAuthService/internal/core"
	"log"
)

type (
	UserService interface {
		Register(username, email, password string) (string, string, error)
		Login(login, password string) (string, string, error)
		FindByID(id uint) (*core.User, error)
	}

	TokenService interface {
		Generate(userID uint) (string, string, error)
		Verify(accessToken string) (userID uint, err error)
		VerifyRefreshToken(refreshToken string) (userID uint, err error)
	}

	AuthenticationService interface {
		Auth(auth string) (*core.User, error)
	}

	Validator interface {
		Email(email string) bool
		Username(username string) bool
		Password(password string) bool
	}

	Handler struct {
		userService           UserService
		tokenService          TokenService
		authenticationService AuthenticationService
		validator             Validator
	}
)

const (
	userLocals = "user"
)

func New(userService UserService, tokenService TokenService, authenticationService AuthenticationService, validator Validator) *Handler {
	return &Handler{
		userService:           userService,
		tokenService:          tokenService,
		authenticationService: authenticationService,
		validator:             validator,
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
