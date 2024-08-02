package service

import "github.com/lapkomo2018/goTwitterAuthService/internal/core"

type RefreshTokenStorage interface {
	Create(refreshToken *core.RefreshToken) error
	First(refreshToken *core.RefreshToken, cond ...interface{}) error
	FindAll(dest interface{}, conds ...interface{}) error
}

type TokenService struct {
	storage RefreshTokenStorage
}

func NewTokenService(refreshTokenStorage RefreshTokenStorage) *TokenService {
	return &TokenService{
		storage: refreshTokenStorage,
	}
}
