package mysql

import (
	"github.com/lapkomo2018/goTwitterAuthService/internal/core"
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
	err := us.db.Model(&core.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (us *UserStorage) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := us.db.Model(&core.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (us *UserStorage) FindByLogin(login string) (*core.User, error) {
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

func (us *UserStorage) FindByID(id int64) (*core.User, error) {
	user := &core.User{}
	err := us.db.Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserStorage) FindByUsername(username string) (*core.User, error) {
	user := &core.User{}
	err := us.db.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserStorage) FindByEmail(email string) (*core.User, error) {
	user := &core.User{}
	err := us.db.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserStorage) Add(user *core.User) error {
	return us.db.Create(user).Error
}

func (us *UserStorage) Save(user *core.User) error {
	return us.db.Save(user).Error
}

func (us *UserStorage) Delete(user *core.User) error {
	return us.db.Delete(user).Error
}
