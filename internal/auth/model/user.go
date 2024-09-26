package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Email       string `gorm:"unique"`
	Username    string `gorm:"unique"`
	Password    string `json:"-"`
	IsActivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

const (
	RoleUser      = 0
	RoleAdminLvl1 = 1
	RoleAdminLvl2 = 2
	RoleAdminLvl3 = 3
)
