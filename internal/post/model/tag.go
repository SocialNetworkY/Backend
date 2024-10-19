package model

type Tag struct {
	ID    uint    `json:"id" gorm:"primary_key"`
	Name  string  `json:"name"`
	Posts []*Post `json:"-" gorm:"many2many:post_tags;"`
}
