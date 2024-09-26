package mysql

import (
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"gorm.io/gorm"
)

type UserStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (us *UserStorage) ExistsByLogin(login string) (bool, error) {
	exists, err := us.ExistsByEmail(login)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	exists, err = us.ExistsByUsername(login)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (us *UserStorage) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := us.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, ErrDatabase
	}
	return count > 0, nil
}

func (us *UserStorage) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := us.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, ErrDatabase
	}
	return count > 0, nil
}

func (us *UserStorage) FindByLogin(login string) (*model.User, error) {
	user, err := us.FindByEmail(login)
	if err == nil {
		return user, nil
	}

	user, err = us.FindByUsername(login)
	if err == nil {
		return user, nil
	}

	return nil, err
}

func (us *UserStorage) Find(id uint) (*model.User, error) {
	user := &model.User{}
	if err := us.db.Where("id = ?", id).First(user).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (us *UserStorage) FindByUsername(username string) (*model.User, error) {
	user := &model.User{}
	if err := us.db.Where("username = ?", username).First(user).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (us *UserStorage) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := us.db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (us *UserStorage) FindByNickname(nickname string) ([]*model.User, error) {
	var users []*model.User
	if err := us.db.Where("nickname = ?", nickname).Find(&users).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return users, nil
}

func (us *UserStorage) Add(user *model.User) error {
	if err := us.db.Create(user).Error; err != nil {
		return ErrUserCreate
	}
	return nil
}

func (us *UserStorage) Save(user *model.User) error {
	if err := us.db.Save(user).Error; err != nil {
		return ErrUserSave
	}
	return nil
}

func (us *UserStorage) Delete(user *model.User) error {
	if err := us.db.Delete(user).Error; err != nil {
		return ErrUserDelete
	}
	return nil
}
