package mysql

import (
	"github.com/lapkomo2018/goTwitterAuthService/internal/core"
	"gorm.io/gorm"
)

type RefreshTokenStorage struct {
	db *gorm.DB
}

func NewRefreshTokenStorage(db *gorm.DB) *RefreshTokenStorage {
	return &RefreshTokenStorage{
		db: db,
	}
}

func (us *RefreshTokenStorage) First(refreshToken *core.RefreshToken, cond ...interface{}) error {
	return us.db.Where(refreshToken).First(refreshToken, cond...).Error
}

func (us *RefreshTokenStorage) FindAll(dest interface{}, conds ...interface{}) error {
	return us.db.Find(dest, conds...).Error
}

func (us *RefreshTokenStorage) Create(refreshToken *core.RefreshToken) error {
	return us.db.Create(refreshToken).Error
}
