package model

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	UserID   uint   `json:"user_id"`
	PostID   uint   `json:"post_id"`
	Content  string `json:"content"`
	Edited   bool   `json:"edited"`
	EditedBy uint   `json:"edited_by"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"edited_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (c *Comment) AfterFind(tx *gorm.DB) (err error) {
	c.Edited = c.EditedBy != 0
	return
}
