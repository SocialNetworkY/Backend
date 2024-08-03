package service

import "fmt"

type (
	RefreshTokenStorage interface {
		SetToken(userID uint, refreshToken string) error
		GetToken(userID uint) (string, error)
	}

	TokenManager interface {
		Generate(userID uint) (accessToken string, refreshToken string, err error)
		Verify(accessToken string) (userID uint, err error)
		VerifyRefreshToken(refreshToken string) (userID uint, err error)
	}
)

type TokenService struct {
	storage      RefreshTokenStorage
	tokenManager TokenManager
}

func NewTokenService(refreshTokenStorage RefreshTokenStorage, tokenManager TokenManager) *TokenService {
	return &TokenService{
		storage:      refreshTokenStorage,
		tokenManager: tokenManager,
	}
}

func (ts *TokenService) Generate(userID uint) (string, string, error) {
	accessToken, refreshToken, err := ts.tokenManager.Generate(userID)
	if err != nil {
		return "", "", err
	}

	if err := ts.storage.SetToken(userID, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (ts *TokenService) Verify(accessToken string) (userID uint, err error) {
	return ts.tokenManager.Verify(accessToken)
}

func (ts *TokenService) VerifyRefreshToken(refreshToken string) (userID uint, err error) {
	userID, err = ts.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return 0, err
	}
	existingToken, err := ts.storage.GetToken(userID)
	if err != nil {
		return 0, err
	}

	if existingToken != refreshToken {
		return 0, fmt.Errorf("invalid refresh token")
	}

	return userID, nil
}
