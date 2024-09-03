package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Email       string `gorm:"unique"`
	Username    string `gorm:"unique"`
	Password    string
	IsActivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
