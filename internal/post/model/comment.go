package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserID   uint   `json:"user_id"`
	PostID   uint   `json:"post_id"`
	Content  string `json:"content"`
	Edited   bool   `json:"edited" gorm:"-"`
	EditedBy uint   `json:"edited_by"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"edited_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (c *Comment) AfterFind(tx *gorm.DB) (err error) {
	c.Edited = c.EditedBy != 0
	return
}

type CommentStatistic struct {
	Total  uint64 `json:"total"`
	Edited uint64 `json:"edited"`
}
