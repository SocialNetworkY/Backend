package core

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique;not null"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	switch {
	case u.Email == "":
		return errors.New("email cannot be empty")
	case u.Username == "":
		return errors.New("username cannot be empty")
	case u.Password == "":
		return errors.New("password cannot be empty")
	default:
		return nil
	}
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	return tx.Where("user_id = ?", u.ID).Delete(&RefreshToken{}).Error
}
