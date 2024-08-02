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

func (us *UserStorage) First(user *core.User, cond ...interface{}) error {
	return us.db.Where(user).First(user, cond...).Error
}

func (us *UserStorage) FindAll(dest interface{}, conds ...interface{}) error {
	return us.db.Find(dest, conds...).Error
}

func (us *UserStorage) Create(user *core.User) error {
	return us.db.Create(user).Error
}
