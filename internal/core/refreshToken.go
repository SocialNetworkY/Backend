package core

import (
	"errors"
	"gorm.io/gorm"
)

type RefreshToken struct {
	UserID uint
	Token  string `gorm:"primaryKey;unique"`
}

func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	switch {
	case rt.UserID == 0:
		return errors.New("user ID cannot be 0")
	case rt.Token == "":
		return errors.New("token cannot be empty")
	default:
		return nil
	}
}
