package mysql

import (
	"errors"
	"github.com/lapkomo2018/goTwitterAuthService/internal/core"
	"gorm.io/gorm"
)

type ActivationTokenStorage struct {
	db *gorm.DB
}

func NewActivationTokenStorage(db *gorm.DB) *ActivationTokenStorage {
	return &ActivationTokenStorage{
		db: db,
	}
}

func (ats *ActivationTokenStorage) Set(userID uint, activationToken string) error {
	var existingToken core.ActivationToken
	err := ats.db.Where("user_id = ?", userID).First(&existingToken).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrActivationTokenSet
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newToken := core.ActivationToken{
			UserID: userID,
			Token:  activationToken,
		}

		if err := ats.db.Create(&newToken).Error; err != nil {
			return ErrActivationTokenCreate
		}
		return nil
	}

	existingToken.Token = activationToken
	if err := ats.db.Save(&existingToken).Error; err != nil {
		return ErrActivationTokenSave
	}
	return nil
}

func (ats *ActivationTokenStorage) Get(userID uint) (*core.ActivationToken, error) {
	existingToken := &core.ActivationToken{}
	err := ats.db.Where("user_id = ?", userID).First(existingToken).Error
	if err != nil {
		return nil, err
	}
	return existingToken, nil
}

func (ats *ActivationTokenStorage) GetByToken(activationToken string) (*core.ActivationToken, error) {
	existingToken := &core.ActivationToken{}
	err := ats.db.Where("token = ?", activationToken).First(existingToken).Error
	if err != nil {
		return nil, ErrActivationTokenNotFound
	}
	return existingToken, nil
}

func (ats *ActivationTokenStorage) Delete(userID uint) error {
	err := ats.db.Where("user_id = ?", userID).Delete(&core.ActivationToken{}).Error
	if err != nil {
		return ErrActivationTokenNotFound
	}
	return nil
}