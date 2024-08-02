package service

import "github.com/lapkomo2018/goTwitterAuthService/internal/core"

type (
	Hasher interface {
		Hash(password string) string
		Verify(hash string, password string) bool
	}

	UserStorage interface {
		First(user *core.User, cond ...interface{}) error
		FindAll(dest interface{}, conds ...interface{}) error
		Create(user *core.User) error
	}

	UserService struct {
		storage      UserStorage
		tokenService *TokenService
		hasher       Hasher
	}
)

// TODO: Implement function Register, Login
func NewUserService(userStorage UserStorage, tokenService *TokenService, hasher Hasher) *UserService {
	return &UserService{
		storage:      userStorage,
		tokenService: tokenService,
		hasher:       hasher,
	}
}

func (us *UserService) Register(username, email, password string) (string, string, error) {
	return "", "", nil
}
