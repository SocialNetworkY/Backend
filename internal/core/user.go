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
	if u.Email == "" {
		err = errors.New("email cannot be empty")
	} else if u.Password == "" {
		err = errors.New("password cannot be empty")
	}
	return
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	tx.Where("user_id = ?", u.ID).Delete(&RefreshToken{})
	return
}
