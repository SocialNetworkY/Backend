package model

type Tag struct {
	ID    uint    `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name"`
	Posts []*Post `json:"-" gorm:"many2many:post_tags;"`
}
