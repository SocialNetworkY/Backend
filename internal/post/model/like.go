package model

import (
	"time"

	"gorm.io/gorm"
)

type Like struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`

	CreatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
