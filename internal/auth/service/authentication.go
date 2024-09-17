package service

import (
	"github.com/lapkomo2018/goTwitterServices/internal/auth/model"
	"strings"
)

type (
	authSchemeHandler func(token string) (*model.User, error)

	AuthenticationUserService interface {
		Find(id uint) (*model.User, error)
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

func (as *AuthenticationService) Auth(auth string) (*model.User, error) {
	authHeaderParts := strings.Split(auth, " ")
	if len(authHeaderParts) != 2 {
		return nil, ErrAuthenticationInvalidAuthString
	}

	scheme := strings.ToLower(authHeaderParts[0])
	tokenString := authHeaderParts[1]

	handler := as.authHandlers[scheme]

	return handler(tokenString)
}

func (as *AuthenticationService) bearerHandler(token string) (*model.User, error) {
	userID, err := as.tokenService.Verify(token)
	if err != nil {
		return nil, err
	}

	user, err := as.userService.Find(userID)
	if err != nil {
		return nil, err
	}

	if !user.IsActivated {
		return nil, ErrUserNotActivated
	}

	return user, nil
}
