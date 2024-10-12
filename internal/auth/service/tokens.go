package service

import "errors"

type (
	TokensRefreshTokenRepo interface {
		Set(userID uint, refreshToken string) error
		Get(userID uint) (string, error)
	}

	TokensManager interface {
		Generate(userID uint) (accessToken string, refreshToken string, err error)
		Verify(accessToken string) (userID uint, err error)
		VerifyRefreshToken(refreshToken string) (userID uint, err error)
	}
)

var (
	ErrRefreshTokenInvalid = errors.New("refresh token is invalid")
)

type TokensService struct {
	repo         TokensRefreshTokenRepo
	tokenManager TokensManager
}

func NewTokensService(repo TokensRefreshTokenRepo, tokenManager TokensManager) *TokensService {
	return &TokensService{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

func (ts *TokensService) Generate(userID uint) (string, string, error) {
	accessToken, refreshToken, err := ts.tokenManager.Generate(userID)
	if err != nil {
		return "", "", err
	}

	if err := ts.repo.Set(userID, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (ts *TokensService) Verify(accessToken string) (userID uint, err error) {
	return ts.tokenManager.Verify(accessToken)
}

func (ts *TokensService) VerifyRefreshToken(refreshToken string) (userID uint, err error) {
	userID, err = ts.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return 0, err
	}
	existingToken, err := ts.repo.Get(userID)
	if err != nil {
		return 0, err
	}

	if existingToken != refreshToken {
		return 0, ErrRefreshTokenInvalid
	}

	return userID, nil
}
