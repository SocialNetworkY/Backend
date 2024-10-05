package mysql

import (
	"errors"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/model"
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

func (us *RefreshTokenStorage) Set(userID uint, refreshToken string) error {
	var existingToken model.RefreshToken
	err := us.db.Where("user_id = ?", userID).First(&existingToken).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newToken := model.RefreshToken{
			UserID: userID,
			Token:  refreshToken,
		}

		if err := us.db.Create(&newToken).Error; err != nil {
			return err
		}
		return nil
	}

	existingToken.Token = refreshToken
	if err := us.db.Save(&existingToken).Error; err != nil {
		return err
	}
	return nil
}

func (us *RefreshTokenStorage) Get(userID uint) (string, error) {
	existingToken := &model.RefreshToken{}
	err := us.db.Where("user_id = ?", userID).First(existingToken).Error
	if err != nil {
		return "", err
	}
	return existingToken.Token, nil
}
