package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/SocialNetworkY/Backend/internal/user/model"
)

type (
	UserRepo interface {
		Add(user *model.User) error
		Save(user *model.User) error
		Delete(user *model.User) error
		ExistsByLogin(login string) (bool, error)
		ExistsByEmail(email string) (bool, error)
		ExistsByUsername(username string) (bool, error)
		FindByLogin(login string) (*model.User, error)
		Find(id uint) (*model.User, error)
		FindSome(skip, limit int) ([]*model.User, error)
		FindByUsername(username string) (*model.User, error)
		FindByEmail(email string) (*model.User, error)
		FindByNickname(nickname string, skip, limit int) ([]*model.User, error)
		Search(query string, skip, limit int) ([]*model.User, error)
		Statistic() (*model.UserStatistic, error)
	}

	AuthGateway interface {
		UpdateUsernameEmail(ctx context.Context, id uint, username, email string) error
		DeleteUser(ctx context.Context, id uint) error
	}

	PostGateway interface {
		DeleteUserPosts(ctx context.Context, id uint) error
		DeleteUserComments(ctx context.Context, id uint) error
		DeleteUserLikes(ctx context.Context, id uint) error
	}

	ReportGateway interface {
		DeleteUserReports(ctx context.Context, userID uint) error
	}

	ImageStorage interface {
		UploadImage(file io.ReadSeeker, filename string) (string, error)
	}

	UserService struct {
		repo UserRepo
		ag   AuthGateway
		pg   PostGateway
		rg   ReportGateway
		is   ImageStorage
	}
)

var (
	ErrUserUsernameTaken = errors.New("username is already taken")
	ErrUserEmailTaken    = errors.New("email is already taken")
)

func NewUserService(repo UserRepo, is ImageStorage, ag AuthGateway, pg PostGateway, rg ReportGateway) *UserService {
	return &UserService{
		repo: repo,
		ag:   ag,
		pg:   pg,
		rg:   rg,
		is:   is,
	}
}

func (us *UserService) Create(id, role uint, username, email string) (*model.User, error) {
	exists, err := us.repo.ExistsByUsername(username)
	switch {
	case err != nil:
		return nil, err
	case exists:
		return nil, ErrUserUsernameTaken
	}

	exists, err = us.repo.ExistsByEmail(email)
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

	if err := us.repo.Add(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Find(id uint) (*model.User, error) {
	return us.repo.Find(id)
}

func (us *UserService) FindByUsername(username string) (*model.User, error) {
	return us.repo.FindByUsername(username)
}

func (us *UserService) FindByEmail(email string) (*model.User, error) {
	return us.repo.FindByEmail(email)
}

func (us *UserService) FindByNickname(nickname string, skip, limit int) ([]*model.User, error) {
	return us.repo.FindByNickname(nickname, skip, limit)
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
		return ErrUserEmailTaken
	}

	user.Email = email
	if err := us.repo.Save(user); err != nil {
		return err
	}

	// Grpc call to update email in auth service
	if err := us.ag.UpdateUsernameEmail(context.Background(), id, "", email); err != nil {
		return err
	}

	return nil
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
		return ErrUserUsernameTaken
	}

	user.Username = username
	if err := us.repo.Save(user); err != nil {
		return err
	}

	// Grpc call to update username in auth service
	if err := us.ag.UpdateUsernameEmail(context.Background(), id, username, ""); err != nil {
		return err
	}

	return nil
}

func (us *UserService) ChangeNickname(id uint, nickname string) error {
	user, err := us.repo.Find(id)
	if err != nil {
		return err
	}

	user.Nickname = nickname
	return us.repo.Save(user)
}

func (us *UserService) ChangeAvatar(id uint, file io.ReadSeeker) error {
	user, err := us.repo.Find(id)
	if err != nil {
		return err
	}

	newAvatarURL, err := us.is.UploadImage(file, fmt.Sprintf("%d_avatar", user.ID))
	if err != nil {
		return err
	}

	user.Avatar = newAvatarURL
	return us.repo.Save(user)
}

func (us *UserService) ChangeRole(id, role uint) error {
	user, err := us.repo.Find(id)
	if err != nil {
		return err
	}

	user.Role = role
	return us.repo.Save(user)
}

func (us *UserService) Delete(id uint) error {
	user, err := us.repo.Find(id)
	if err != nil {
		return err
	}

	if err := us.pg.DeleteUserPosts(context.Background(), id); err != nil {
		return err
	}

	if err := us.pg.DeleteUserComments(context.Background(), id); err != nil {
		return err
	}

	if err := us.pg.DeleteUserLikes(context.Background(), id); err != nil {
		return err
	}

	if err := us.rg.DeleteUserReports(context.Background(), id); err != nil {
		return err
	}

	if err := us.ag.DeleteUser(context.Background(), id); err != nil {
		return err
	}

	return us.repo.Delete(user)
}

func (us *UserService) Search(query string, skip, limit int) ([]*model.User, error) {
	return us.repo.Search(query, skip, limit)
}

func (us *UserService) Statistic() (*model.UserStatistic, error) {
	return us.repo.Statistic()
}
