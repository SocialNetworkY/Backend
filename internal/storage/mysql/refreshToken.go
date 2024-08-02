package mysql

import (
	"errors"
	"github.com/lapkomo2018/goTwitterAuthService/internal/core"
	"gorm.io/gorm"
)

type (
	RefreshTokenStorage struct {
		db *gorm.DB
	}
)

func NewRefreshTokenStorage(db *gorm.DB) *RefreshTokenStorage {
	return &RefreshTokenStorage{
		db: db,
	}
}

func (us *RefreshTokenStorage) SetToken(userID uint, refreshToken string) error {
	var existingToken core.RefreshToken
	err := us.db.Where("user_id = ?", userID).First(&existingToken).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newToken := core.RefreshToken{
			UserID: userID,
			Token:  refreshToken,
		}

		return us.db.Create(&newToken).Error
	}

	existingToken.Token = refreshToken
	return us.db.Save(&existingToken).Error
}
