package repository

import (
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (pr *PostRepository) Add(post *model.Post) error {
	if err := pr.db.Create(post).Error; err != nil {
		return err
	}
	return nil
}

func (pr *PostRepository) Save(post *model.Post) error {
	if err := pr.db.Save(post).Error; err != nil {
		return err
	}
	return nil
}

func (pr *PostRepository) Delete(id uint) error {
	if err := pr.db.Delete(&model.Post{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

// Find finds a post by id
func (pr *PostRepository) Find(id uint) (*model.Post, error) {
	post := &model.Post{}
	if err := pr.db.Preload("Tags").Preload("Comments").Preload("Likes").First(post, id).Error; err != nil {
		return nil, err
	}
	return post, nil
}

// FindSome fetches some posts
func (pr *PostRepository) FindSome(skip, limit int) ([]*model.Post, error) {
	var posts []*model.Post
	if err := pr.db.Preload("Tags").Preload("Comments").Preload("Likes").Offset(skip).Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindByUser fetches some posts by user id
func (pr *PostRepository) FindByUser(userID uint, skip, limit int) ([]*model.Post, error) {
	var posts []*model.Post
	if err := pr.db.Preload("Tags").Preload("Comments").Preload("Likes").Offset(skip).Limit(limit).Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindByTag fetches some posts by tag id
func (pr *PostRepository) FindByTag(tagID uint, skip, limit int) ([]*model.Post, error) {
	var posts []*model.Post
	if err := pr.db.Preload("Tags").Preload("Comments").Preload("Likes").Offset(skip).Limit(limit).Joins("JOIN post_tags ON post_tags.post_id = posts.id").Where("post_tags.tag_id = ?", tagID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
