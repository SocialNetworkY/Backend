package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterAuthService/internal/core"
	"log"
)

type (
	UserService interface {
		Register(username, email, password string) (string, error)
		Login(login, password string) (string, string, error)
		Activate(token string) (string, string, error)
		Find(id uint) (*core.User, error)
		ChangeEmail(id uint, email string) error
		ChangeUsername(id uint, username string) error
		ChangePassword(id uint, password string) error
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
		Login(login string) error
		Email(email string) error
		Username(username string) error
		Password(password string) error
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
		h.initUserChangeApi(v1)
	}
}
