package model

type Tag struct {
	ID    uint    `json:"id" gorm:"primary_key"`
	Name  string  `json:"name" gorm:"unique"`
	Posts []*Post `json:"-" gorm:"many2many:post_tags;"`
}
