package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"gorm.io/gorm"
)

// StringArray is a custom type for []string
type StringArray []string

// Value implements the driver.Valuer interface
func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface
func (a *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, a)
}

type Post struct {
	ID             uint        `json:"id" gorm:"primaryKey"`
	UserID         uint        `json:"user_id"`
	Title          string      `json:"title"`
	Content        string      `json:"content"`
	ImageURLs      StringArray `json:"image_urls"`
	VideoURLs      StringArray `json:"video_urls"`
	Tags           []*Tag      `json:"tags" gorm:"many2many:post_tags;"`
	TagsAmount     uint        `json:"tags_amount" gorm:"-"`
	Comments       []*Comment  `json:"comments"`
	CommentsAmount uint        `json:"comments_amount" gorm:"-"`
	Likes          []*Like     `json:"likes"`
	LikesAmount    uint        `json:"likes_amount" gorm:"-"`
	Liked          bool        `json:"is_liked" gorm:"-"`
	PostedAt       time.Time   `json:"posted_at"`
	Edited         bool        `json:"edited" gorm:"-"`
	EditedBy       uint        `json:"edited_by"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"edited_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// AfterFind loads all tags, comments and likes associated with the post
func (p *Post) AfterFind(tx *gorm.DB) (err error) {
	p.TagsAmount = uint(len(p.Tags))
	p.CommentsAmount = uint(len(p.Comments))
	p.LikesAmount = uint(len(p.Likes))
	for _, like := range p.Likes {
		if p.UserID == like.UserID {
			p.Liked = true
			break
		}
	}
	sort.Slice(p.Comments, func(i, j int) bool {
		return p.Comments[i].CreatedAt.After(p.Comments[j].CreatedAt)
	})
	p.Edited = p.EditedBy != 0
	return nil
}
