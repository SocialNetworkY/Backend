package mysql

import (
	"errors"
)

var (
	ErrDatabase  = errors.New("database error")
	ErrBDConnect = errors.New("error connecting to the database")
	ErrBDPing    = errors.New("error pinging the database")

	ErrBanCreate = errors.New("error creating ban")
	ErrBanSave   = errors.New("error saving ban")
	ErrBanDelete = errors.New("error deleting ban")
	ErrBanFind   = errors.New("error finding ban")

	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
	ErrUserCreate   = errors.New("error creating user")
	ErrUserSave     = errors.New("error saving user")
	ErrUserDelete   = errors.New("error deleting user")
)
