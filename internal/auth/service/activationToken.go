package service

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/lapkomo2018/goTwitterServices/pkg/model"
)

type (
	ActivationTokenStorage interface {
		Set(userID uint, activationToken string) error
		Get(userID uint) (*model.ActivationToken, error)
		GetByToken(activationToken string) (*model.ActivationToken, error)
		Delete(userID uint) error
	}
)

type ActivationTokenService struct {
	storage ActivationTokenStorage
}

func NewActivationTokenService(storage ActivationTokenStorage) *ActivationTokenService {
	return &ActivationTokenService{
		storage: storage,
	}
}

func (ts *ActivationTokenService) Generate(userID uint) (string, error) {
	// generate activation token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", ErrActivationTokenGeneration
	}
	activationToken := hex.EncodeToString(tokenBytes)

	if err := ts.storage.Set(userID, activationToken); err != nil {
		return "", err
	}

	return activationToken, nil
}

func (ts *ActivationTokenService) Get(userID uint) (*model.ActivationToken, error) {
	return ts.storage.Get(userID)
}

func (ts *ActivationTokenService) GetByToken(activationToken string) (*model.ActivationToken, error) {
	return ts.storage.GetByToken(activationToken)
}

func (ts *ActivationTokenService) Delete(userID uint) error {
	return ts.storage.Delete(userID)
}
