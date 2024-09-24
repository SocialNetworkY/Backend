package mysql

import (
	"errors"
)

var (
	ErrDatabase  = errors.New("database error")
	ErrBDConnect = errors.New("error connecting to the database")
	ErrBDPing    = errors.New("error pinging the database")

	ErrActivationTokenSet      = errors.New("error setting activation token")
	ErrActivationTokenCreate   = errors.New("error creating activation token")
	ErrActivationTokenSave     = errors.New("error saving activation token")
	ErrActivationTokenNotFound = errors.New("activation token not found")

	ErrRefreshTokenSet      = errors.New("error setting refresh token")
	ErrRefreshTokenCreate   = errors.New("error creating refresh token")
	ErrRefreshTokenSave     = errors.New("error saving refresh token")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")

	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
	ErrUserCreate   = errors.New("error creating user")
	ErrUserSave     = errors.New("error saving user")
	ErrUserDelete   = errors.New("error deleting user")
)
