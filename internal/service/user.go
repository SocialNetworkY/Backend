package service

import (
	"errors"
	"github.com/lapkomo2018/goTwitterAuthService/internal/core"
)

type (
	Hasher interface {
		Hash(password string) string
		Verify(hash string, password string) bool
	}

	UserStorage interface {
		ExistsByLogin(login string) (bool, error)
		ExistsByEmail(email string) (bool, error)
		ExistsByUsername(username string) (bool, error)
		FindByLogin(login string) (*core.User, error)
		FindByID(id uint) (*core.User, error)
		FindByUsername(username string) (*core.User, error)
		FindByEmail(email string) (*core.User, error)
		Add(user *core.User) error
		Save(user *core.User) error
		Delete(user *core.User) error
	}

	UserTokenService interface {
		Generate(userID uint) (string, string, error)
	}

	UserService struct {
		storage      UserStorage
		tokenService UserTokenService
		hasher       Hasher
	}
)

func NewUserService(userStorage UserStorage, tokenService UserTokenService, hasher Hasher) *UserService {
	return &UserService{
		storage:      userStorage,
		tokenService: tokenService,
		hasher:       hasher,
	}
}

func (us *UserService) Register(username, email, password string) (string, string, error) {
	exists, err := us.storage.ExistsByUsername(username)
	switch {
	case err != nil:
		return "", "", err
	case exists:
		return "", "", errors.New("username already taken")
	}

	exists, err = us.storage.ExistsByEmail(email)
	switch {
	case err != nil:
		return "", "", err
	case exists:
		return "", "", errors.New("email already taken")
	}

	hashedPassword := us.hasher.Hash(password)

	user := &core.User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
	}

	if err := us.storage.Add(user); err != nil {
		return "", "", err
	}

	return us.tokenService.Generate(user.ID)
}

func (us *UserService) Login(login, password string) (string, string, error) {
	user, err := us.storage.FindByLogin(login)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	if !us.hasher.Verify(user.Password, password) {
		return "", "", errors.New("invalid password")
	}

	return us.tokenService.Generate(user.ID)
}

func (us *UserService) FindByID(id uint) (*core.User, error) {
	return us.storage.FindByID(id)
}
