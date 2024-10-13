package service

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/SocialNetworkY/Backend/internal/auth/model"
)

type (
	ActivationTokenRepo interface {
		Set(userID uint, activationToken string) error
		Get(userID uint) (*model.ActivationToken, error)
		GetByToken(activationToken string) (*model.ActivationToken, error)
		Delete(userID uint) error
	}
)

type ActivationTokenService struct {
	repo ActivationTokenRepo
}

func NewActivationTokenService(repo ActivationTokenRepo) *ActivationTokenService {
	return &ActivationTokenService{
		repo: repo,
	}
}

func (ts *ActivationTokenService) Generate(userID uint) (string, error) {
	// generate activation token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	activationToken := hex.EncodeToString(tokenBytes)

	if err := ts.repo.Set(userID, activationToken); err != nil {
		return "", err
	}

	return activationToken, nil
}

func (ts *ActivationTokenService) Get(userID uint) (*model.ActivationToken, error) {
	return ts.repo.Get(userID)
}

func (ts *ActivationTokenService) GetByToken(activationToken string) (*model.ActivationToken, error) {
	return ts.repo.GetByToken(activationToken)
}

func (ts *ActivationTokenService) Delete(userID uint) error {
	return ts.repo.Delete(userID)
}
