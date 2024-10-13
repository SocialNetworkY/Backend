package repository

import (
	"errors"
	"github.com/SocialNetworkY/Backend/internal/auth/model"
	"gorm.io/gorm"
)

type (
	RefreshTokenRepository struct {
		db *gorm.DB
	}
)

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

func (rtr *RefreshTokenRepository) Set(userID uint, refreshToken string) error {
	var existingToken model.RefreshToken
	err := rtr.db.Where("user_id = ?", userID).First(&existingToken).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newToken := model.RefreshToken{
			UserID: userID,
			Token:  refreshToken,
		}

		if err := rtr.db.Create(&newToken).Error; err != nil {
			return err
		}
		return nil
	}

	existingToken.Token = refreshToken
	if err := rtr.db.Save(&existingToken).Error; err != nil {
		return err
	}
	return nil
}

func (rtr *RefreshTokenRepository) Get(userID uint) (string, error) {
	existingToken := &model.RefreshToken{}
	err := rtr.db.Where("user_id = ?", userID).First(existingToken).Error
	if err != nil {
		return "", err
	}
	return existingToken.Token, nil
}
