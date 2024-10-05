package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"io"
)

type (
	UserStorage interface {
		ExistsByLogin(login string) (bool, error)
		ExistsByEmail(email string) (bool, error)
		ExistsByUsername(username string) (bool, error)
		FindByLogin(login string) (*model.User, error)
		Find(id uint) (*model.User, error)
		FindByUsername(username string) (*model.User, error)
		FindByEmail(email string) (*model.User, error)
		FindByNickname(nickname string) ([]*model.User, error)
		Add(user *model.User) error
		Save(user *model.User) error
		Delete(user *model.User) error
	}

	AuthGateway interface {
		UpdateUsernameEmail(ctx context.Context, auth string, id uint, username, email string) error
		DeleteUser(ctx context.Context, auth string, id uint) error
	}

	ImageStorage interface {
		UploadImage(file io.ReadSeeker, filename string) (string, error)
	}

	UserService struct {
		us UserStorage
		ag AuthGateway
		is ImageStorage
	}
)

var (
	ErrUserUsernameTaken = errors.New("username is already taken")
	ErrUserEmailTaken    = errors.New("email is already taken")
)

func NewUserService(us UserStorage, is ImageStorage, ag AuthGateway) *UserService {
	return &UserService{
		us: us,
		ag: ag,
		is: is,
	}
}

func (us *UserService) Create(id, role uint, username, email string) (*model.User, error) {
	exists, err := us.us.ExistsByUsername(username)
	switch {
	case err != nil:
		return nil, err
	case exists:
		return nil, ErrUserUsernameTaken
	}

	exists, err = us.us.ExistsByEmail(email)
	switch {
	case err != nil:
		return nil, err
	case exists:
		return nil, ErrUserEmailTaken
	}

	user := &model.User{
		ID:       id,
		Email:    email,
		Username: username,
		Nickname: username,
		Role:     role,
	}

	if err := us.us.Add(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Find(id uint) (*model.User, error) {
	return us.us.Find(id)
}

func (us *UserService) FindByUsername(username string) (*model.User, error) {
	return us.us.FindByUsername(username)
}

func (us *UserService) FindByEmail(email string) (*model.User, error) {
	return us.us.FindByEmail(email)
}

func (us *UserService) FindByNickname(nickname string) ([]*model.User, error) {
	return us.us.FindByNickname(nickname)
}

func (us *UserService) ChangeEmail(id uint, auth, email string) error {
	user, err := us.us.Find(id)
	if err != nil {
		return err
	}

	exists, err := us.us.ExistsByEmail(email)
	switch {
	case err != nil:
		return err
	case exists:
		return ErrUserEmailTaken
	}

	user.Email = email
	if err := us.us.Save(user); err != nil {
		return err
	}

	// Grpc call to update email in auth service
	if err := us.ag.UpdateUsernameEmail(context.Background(), auth, id, "", email); err != nil {
		return err
	}

	return nil
}

func (us *UserService) ChangeUsername(id uint, auth, username string) error {
	user, err := us.us.Find(id)
	if err != nil {
		return err
	}

	exists, err := us.us.ExistsByUsername(username)
	switch {
	case err != nil:
		return err
	case exists:
		return ErrUserUsernameTaken
	}

	user.Username = username
	if err := us.us.Save(user); err != nil {
		return err
	}

	// Grpc call to update username in auth service
	if err := us.ag.UpdateUsernameEmail(context.Background(), auth, id, username, ""); err != nil {
		return err
	}

	return nil
}

func (us *UserService) ChangeNickname(id uint, nickname string) error {
	user, err := us.us.Find(id)
	if err != nil {
		return err
	}

	user.Nickname = nickname
	return us.us.Save(user)
}

func (us *UserService) ChangeAvatar(id uint, file io.ReadSeeker) error {
	user, err := us.us.Find(id)
	if err != nil {
		return err
	}

	newAvatarURL, err := us.is.UploadImage(file, fmt.Sprintf("%d_avatar", user.ID))
	if err != nil {
		return err
	}

	user.Avatar = newAvatarURL
	return us.us.Save(user)
}

func (us *UserService) ChangeRole(id, role uint) error {
	user, err := us.us.Find(id)
	if err != nil {
		return err
	}

	user.Role = role
	return us.us.Save(user)
}

func (us *UserService) Delete(id uint, auth string) error {
	user, err := us.us.Find(id)
	if err != nil {
		return err
	}

	// Grpc call to delete user from other services
	if err := us.ag.DeleteUser(context.Background(), auth, id); err != nil {
		return err
	}

	if err := us.us.Delete(user); err != nil {
		return err
	}

	return nil
}
