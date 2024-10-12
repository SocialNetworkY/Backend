package model

import (
	"gorm.io/gorm"
)

type Tag struct {
	ID    uint    `json:"id" gorm:"primary_key"`
	Name  string  `json:"name" gorm:"unique"`
	Posts []*Post `json:"-" gorm:"many2many:post_tags;"`
}

// AfterDelete removes all references to the tag from posts
func (t *Tag) AfterDelete(tx *gorm.DB) (err error) {
	if err = tx.Association("Posts").Clear(); err != nil {
		return err
	}

	return nil
}
