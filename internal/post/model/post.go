package model

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID             uint       `json:"id" gorm:"primary_key"`
	UserID         uint       `json:"user_id"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	ImageURLs      []string   `json:"image_urls"`
	VideoURLs      []string   `json:"video_urls"`
	Tags           []*Tag     `json:"tags" gorm:"many2many:post_tags;"`
	TagsAmount     uint       `json:"tags_amount" gorm:"-"`
	Comments       []*Comment `json:"comments"`
	CommentsAmount uint       `json:"comments_amount" gorm:"-"`
	Likes          []*Like    `json:"likes"`
	LikesAmount    uint       `json:"likes_amount" gorm:"-"`
	PostedAt       time.Time  `json:"posted_at"`
	Updated        bool       `json:"updated" gorm:"-"`
	UpdatedBy      uint       `json:"updated_by"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// AfterFind loads all tags, comments and likes associated with the post
func (p *Post) AfterFind(tx *gorm.DB) (err error) {
	if err = tx.Model(p).Association("Tags").Find(&p.Tags); err != nil {
		return err
	}
	p.TagsAmount = uint(len(p.Tags))

	if err = tx.Model(p).Association("Comments").Find(&p.Comments); err != nil {
		return err
	}
	p.CommentsAmount = uint(len(p.Comments))

	if err = tx.Model(p).Association("Likes").Find(&p.Likes); err != nil {
		return err
	}
	p.LikesAmount = uint(len(p.Likes))

	p.Updated = p.UpdatedAt.After(p.PostedAt)

	return nil
}

// AfterDelete removes all references to the post from tags, deletes posts comments and likes
func (p *Post) AfterDelete(tx *gorm.DB) (err error) {
	if err = tx.Model(p).Association("Tags").Clear(); err != nil {
		return err
	}

	if err = tx.Where("post_id = ?", p.ID).Delete(&Comment{}).Error; err != nil {
		return err
	}

	if err = tx.Where("post_id = ?", p.ID).Delete(&Like{}).Error; err != nil {
		return err
	}

	return nil
}
