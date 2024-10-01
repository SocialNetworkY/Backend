package service

import (
	"context"
	"fmt"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/model"
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
		FindByLogin(login string) (*model.User, error)
		Find(id uint) (*model.User, error)
		FindByUsername(username string) (*model.User, error)
		FindByEmail(email string) (*model.User, error)
		Add(user *model.User) error
		Save(user *model.User) error
		Delete(user *model.User) error
	}

	UserTokenService interface {
		Generate(userID uint) (string, string, error)
	}

	UserActivationTokenService interface {
		Generate(userID uint) (string, error)
		Get(userID uint) (*model.ActivationToken, error)
		GetByToken(activationToken string) (*model.ActivationToken, error)
		Delete(userID uint) error
	}

	UserGateway interface {
		CreateUser(ctx context.Context, auth string, userID, role uint, username, email string) error
		GetUserRole(ctx context.Context, auth string, userID uint) (uint, error)
	}

	UserService struct {
		storage                UserStorage
		tokenService           UserTokenService
		activationTokenService UserActivationTokenService
		hasher                 Hasher
		userGateway            UserGateway
	}
)

func NewUserService(userStorage UserStorage, tokenService UserTokenService, activationTokenService UserActivationTokenService, hasher Hasher, ug UserGateway) *UserService {
	return &UserService{
		storage:                userStorage,
		tokenService:           tokenService,
		activationTokenService: activationTokenService,
		hasher:                 hasher,
		userGateway:            ug,
	}
}

func (us *UserService) Register(username, email, password string) (activationToken string, err error) {
	exists, err := us.storage.ExistsByUsername(username)
	switch {
	case err != nil:
		return "", err
	case exists:
		return "", ErrUserUsernameTaken
	}

	exists, err = us.storage.ExistsByEmail(email)
	switch {
	case err != nil:
		return "", err
	case exists:
		return "", ErrUserEmailTaken
	}

	hashedPassword := us.hasher.Hash(password)

	user := &model.User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
	}

	if err := us.storage.Add(user); err != nil {
		return "", err
	}

	activationToken, err = us.activationTokenService.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return activationToken, nil
}

func (us *UserService) Login(login, password string) (string, string, error) {
	user, err := us.storage.FindByLogin(login)
	if err != nil {
		return "", "", err
	}

	if !us.hasher.Verify(user.Password, password) {
		return "", "", ErrUserInvalidPassword
	}

	if !user.IsActivated {
		return "", "", ErrUserNotActivated
	}

	return us.tokenService.Generate(user.ID)
}

func (us *UserService) Activate(activationToken string) (accessToken, refreshToken string, err error) {
	token, err := us.activationTokenService.GetByToken(activationToken)
	if err != nil {
		return "", "", err
	}

	user, err := us.storage.Find(token.UserID)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err = us.tokenService.Generate(user.ID)
	if err != nil {
		return "", "", err
	}

	if err := us.userGateway.CreateUser(context.Background(), fmt.Sprintf("Bearer %s", accessToken), user.ID, model.RoleUser, user.Username, user.Email); err != nil {
		return "", "", err
	}

	user.IsActivated = true
	if err := us.storage.Save(user); err != nil {
		return "", "", err
	}

	if err := us.activationTokenService.Delete(token.UserID); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (us *UserService) Find(id uint) (*model.User, error) {
	return us.storage.Find(id)
}

func (us *UserService) ChangeEmail(id uint, email string) error {
	user, err := us.storage.Find(id)
	if err != nil {
		return err
	}

	exists, err := us.storage.ExistsByEmail(email)
	switch {
	case err != nil:
		return err
	case exists:
		return ErrUserEmailTaken
	}

	user.Email = email
	return us.storage.Save(user)
}

func (us *UserService) ChangeUsername(id uint, username string) error {
	user, err := us.storage.Find(id)
	if err != nil {
		return err
	}

	exists, err := us.storage.ExistsByUsername(username)
	switch {
	case err != nil:
		return err
	case exists:
		return ErrUserUsernameTaken
	}

	user.Username = username
	return us.storage.Save(user)
}

func (us *UserService) ChangePassword(id uint, password string) error {
	user, err := us.storage.Find(id)
	if err != nil {
		return err
	}

	hashedPassword := us.hasher.Hash(password)
	user.Password = hashedPassword
	return us.storage.Save(user)
}

func (us *UserService) Delete(id uint) error {
	user, err := us.storage.Find(id)
	if err != nil {
		return err
	}

	return us.storage.Delete(user)
}
