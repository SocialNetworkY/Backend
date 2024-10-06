package model

import "gorm.io/gorm"

type Tag struct {
	ID    uint    `json:"id" gorm:"primary_key"`
	Name  string  `json:"name" gorm:"unique"`
	Posts []*Post `json:"posts" gorm:"many2many:post_tags;"`
}

// AfterFind loads all posts associated with the tag
func (t *Tag) AfterFind(tx *gorm.DB) (err error) {
	if err = tx.Model(t).Association("Posts").Find(&t.Posts); err != nil {
		return err
	}

	return nil
}

// AfterDelete removes all references to the tag from posts
func (t *Tag) AfterDelete(tx *gorm.DB) (err error) {
	if err = tx.Model(t).Association("Posts").Clear(); err != nil {
		return err
	}

	return nil
}
