package repository

import (
	"errors"
	"github.com/SocialNetworkY/Backend/internal/auth/model"
	"gorm.io/gorm"
)

type ActivationTokenRepository struct {
	db *gorm.DB
}

func NewActivationTokenRepository(db *gorm.DB) *ActivationTokenRepository {
	return &ActivationTokenRepository{
		db: db,
	}
}

func (atr *ActivationTokenRepository) Set(userID uint, activationToken string) error {
	var existingToken model.ActivationToken
	err := atr.db.Where("user_id = ?", userID).First(&existingToken).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newToken := model.ActivationToken{
			UserID: userID,
			Token:  activationToken,
		}

		if err := atr.db.Create(&newToken).Error; err != nil {
			return err
		}
		return nil
	}

	existingToken.Token = activationToken
	if err := atr.db.Save(&existingToken).Error; err != nil {
		return err
	}
	return nil
}

func (atr *ActivationTokenRepository) Get(userID uint) (*model.ActivationToken, error) {
	existingToken := &model.ActivationToken{}
	err := atr.db.Where("user_id = ?", userID).First(existingToken).Error
	if err != nil {
		return nil, err
	}
	return existingToken, nil
}

func (atr *ActivationTokenRepository) GetByToken(activationToken string) (*model.ActivationToken, error) {
	existingToken := &model.ActivationToken{}
	err := atr.db.Where("token = ?", activationToken).First(existingToken).Error
	if err != nil {
		return nil, err
	}
	return existingToken, nil
}

func (atr *ActivationTokenRepository) Delete(userID uint) error {
	err := atr.db.Where("user_id = ?", userID).Delete(&model.ActivationToken{}).Error
	if err != nil {
		return err
	}
	return nil
}
