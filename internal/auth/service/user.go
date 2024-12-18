package service

import (
	"context"
	"errors"

	"github.com/SocialNetworkY/Backend/internal/auth/model"
	"github.com/SocialNetworkY/Backend/pkg/constant"
)

type (
	Hasher interface {
		Hash(password string) string
		Verify(hash string, password string) bool
	}

	UserRepo interface {
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
		CreateUser(ctx context.Context, userID, role uint, username, email string) error
	}

	UserService struct {
		repo                   UserRepo
		tokenService           UserTokenService
		activationTokenService UserActivationTokenService
		hasher                 Hasher
		userGateway            UserGateway
	}
)

var (
	ErrUsernameTaken    = errors.New("username is taken")
	ErrEmailTaken       = errors.New("email is taken")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUserNotActivated = errors.New("user is not activated")
)

func NewUserService(repo UserRepo, tokenService UserTokenService, activationTokenService UserActivationTokenService, hasher Hasher, ug UserGateway) *UserService {
	return &UserService{
		repo:                   repo,
		tokenService:           tokenService,
		activationTokenService: activationTokenService,
		hasher:                 hasher,
		userGateway:            ug,
	}
}

func (us *UserService) Register(username, email, password string) (activationToken string, err error) {
	exists, err := us.repo.ExistsByUsername(username)
	switch {
	case err != nil:
		return "", err
	case exists:
		return "", ErrUsernameTaken
	}

	exists, err = us.repo.ExistsByEmail(email)
	switch {
	case err != nil:
		return "", err
	case exists:
		return "", ErrEmailTaken
	}

	hashedPassword := us.hasher.Hash(password)

	user := &model.User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
	}

	if err := us.repo.Add(user); err != nil {
		return "", err
	}

	activationToken, err = us.activationTokenService.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return activationToken, nil
}

func (us *UserService) Login(login, password string) (string, string, error) {
	user, err := us.repo.FindByLogin(login)
	if err != nil {
		return "", "", err
	}

	if !us.hasher.Verify(user.Password, password) {
		return "", "", ErrInvalidPassword
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

	user, err := us.repo.Find(token.UserID)
	if err != nil {
		return "", "", err
	}

	if err := us.userGateway.CreateUser(context.Background(), user.ID, constant.RoleUser, user.Username, user.Email); err != nil {
		return "", "", err
	}

	user.IsActivated = true
	if err := us.repo.Save(user); err != nil {
		return "", "", err
	}

	if err := us.activationTokenService.Delete(token.UserID); err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err = us.tokenService.Generate(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (us *UserService) Find(id uint) (*model.User, error) {
	return us.repo.Find(id)
}

func (us *UserService) ChangeEmail(id uint, email string) error {
	user, err := us.repo.Find(id)
	if err != nil {
		return err
	}

	exists, err := us.repo.ExistsByEmail(email)
	switch {
	case err != nil:
		return err
	case exists:
		return ErrEmailTaken
	}

	user.Email = email
	return us.repo.Save(user)
}

func (us *UserService) ChangeUsername(id uint, username string) error {
	user, err := us.repo.Find(id)
	if err != nil {
		return err
	}

	exists, err := us.repo.ExistsByUsername(username)
	switch {
	case err != nil:
		return err
	case exists:
		return ErrUsernameTaken
	}

	user.Username = username
	return us.repo.Save(user)
}

func (us *UserService) ChangePassword(id uint, password string) error {
	user, err := us.repo.Find(id)
	if err != nil {
		return err
	}

	hashedPassword := us.hasher.Hash(password)
	user.Password = hashedPassword
	return us.repo.Save(user)
}

func (us *UserService) Delete(id uint) error {
	user, err := us.repo.Find(id)
	if err != nil {
		return err
	}

	return us.repo.Delete(user)
}
