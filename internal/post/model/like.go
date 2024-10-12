package model

import (
	"gorm.io/gorm"
	"time"
)

type Like struct {
	ID     uint `json:"id" gorm:"primary_key"`
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`

	CreatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
