package core

import "errors"

var (
	ErrUserLoadingFailed         = errors.New("user loading failed")
	ErrUserRefreshTokenDelete    = errors.New("user refresh token delete failed")
	ErrUserActivationTokenDelete = errors.New("user activation token delete failed")
)
