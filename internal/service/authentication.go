package service

import (
	"errors"
	"github.com/lapkomo2018/goTwitterAuthService/internal/core"
	"strings"
)

type (
	authSchemeHandler func(token string) (*core.User, error)

	AuthenticationUserService interface {
		FindByID(id uint) (*core.User, error)
	}
	AuthenticationTokenService interface {
		Verify(accessToken string) (userID uint, err error)
	}

	AuthenticationService struct {
		userService  AuthenticationUserService
		tokenService AuthenticationTokenService
		authHandlers map[string]authSchemeHandler
	}
)

func NewAuthenticationService(us AuthenticationUserService, ts AuthenticationTokenService) *AuthenticationService {
	authService := &AuthenticationService{
		userService:  us,
		tokenService: ts,
	}

	authService.authHandlers = map[string]authSchemeHandler{
		"bearer": authService.bearerHandler,
	}

	return authService
}

func (as *AuthenticationService) Auth(auth string) (*core.User, error) {
	authHeaderParts := strings.Split(auth, " ")
	if len(authHeaderParts) != 2 {
		return nil, errors.New("invalid auth string")
	}

	scheme := strings.ToLower(authHeaderParts[0])
	tokenString := authHeaderParts[1]

	handler := as.authHandlers[scheme]

	return handler(tokenString)
}

func (as *AuthenticationService) bearerHandler(token string) (*core.User, error) {
	userID, err := as.tokenService.Verify(token)
	if err != nil {
		return nil, err
	}

	return as.userService.FindByID(userID)
}
