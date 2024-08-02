package service

import "github.com/lapkomo2018/goTwitterAuthService/internal/core"

type UserStorage interface {
	First(user *core.User, cond ...interface{}) error
	FindAll(dest interface{}, conds ...interface{}) error
	Create(user *core.User) error
}

type UserService struct {
	storage      UserStorage
	tokenService *TokenService
}

func NewUserService(userStorage UserStorage, tokenService *TokenService) *UserService {
	return &UserService{
		storage:      userStorage,
		tokenService: tokenService,
	}
}
