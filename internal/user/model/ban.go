package model

import (
	"gorm.io/gorm"
	"time"
)

type Ban struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	AdminID   uint
	Reason    string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
