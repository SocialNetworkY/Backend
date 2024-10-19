package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint `gorm:"primaryKey"`
	Email       string
	Username    string
	Password    string `json:"-"`
	IsActivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
