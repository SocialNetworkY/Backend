package service

import "errors"

var (
	ErrActivationTokenGeneration = errors.New("error generating activation token")

	ErrAuthenticationInvalidAuthString = errors.New("invalid auth string")

	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	ErrUserNotActivated    = errors.New("user not activated")
	ErrUserUsernameTaken   = errors.New("username already taken")
	ErrUserEmailTaken      = errors.New("email already taken")
	ErrUserInvalidPassword = errors.New("invalid password")
)
