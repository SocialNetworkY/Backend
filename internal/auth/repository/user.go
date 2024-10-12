package repository

import (
	"github.com/lapkomo2018/goTwitterServices/internal/auth/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) ExistsByLogin(login string) (bool, error) {
	exists, err := ur.ExistsByEmail(login)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	exists, err = ur.ExistsByUsername(login)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (ur *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := ur.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ur *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := ur.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ur *UserRepository) FindByLogin(login string) (*model.User, error) {
	user, err := ur.FindByEmail(login)
	if err == nil {
		return user, nil
	}

	user, err = ur.FindByUsername(login)
	if err == nil {
		return user, nil
	}

	return nil, err
}

func (ur *UserRepository) Find(id uint) (*model.User, error) {
	user := &model.User{}
	if err := ur.db.Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) FindByUsername(username string) (*model.User, error) {
	user := &model.User{}
	if err := ur.db.Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := ur.db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) Add(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) Save(user *model.User) error {
	if err := ur.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) Delete(user *model.User) error {
	if err := ur.db.Delete(user).Error; err != nil {
		return err
	}
	return nil
}
