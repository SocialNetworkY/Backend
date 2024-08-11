package jwt

import "errors"

var (
	ErrGenerateAccessToken     = errors.New("error generating access token")
	ErrGenerateRefreshToken    = errors.New("error generating refresh token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)
