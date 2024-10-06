package model

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	UserID      uint      `json:"user_id"`
	PostID      uint      `json:"post_id"`
	Content     string    `json:"content"`
	CommentedAt time.Time `json:"commented_at"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
