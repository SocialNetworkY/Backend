package mysql

import (
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
	"gorm.io/gorm"
)

type PostStorage struct {
	db *gorm.DB
}

func NewPostStorage(db *gorm.DB) *PostStorage {
	return &PostStorage{
		db: db,
	}
}

func (us *PostStorage) Add(post *model.Post) error {
	if err := us.db.Create(post).Error; err != nil {
		return err
	}
	return nil
}

func (us *PostStorage) Save(post *model.Post) error {
	if err := us.db.Save(post).Error; err != nil {
		return err
	}
	return nil
}

func (us *PostStorage) Delete(id uint) error {
	if err := us.db.Delete(&model.Post{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Find finds a post by id
func (us *PostStorage) Find(id uint) (*model.Post, error) {
	post := &model.Post{}
	if err := us.db.First(post, id).Error; err != nil {
		return nil, err
	}
	return post, nil
}

// FindAll finds all posts
func (us *PostStorage) FindAll() ([]*model.Post, error) {
	var posts []*model.Post
	if err := us.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindFrom fetches posts from 'from'
func (us *PostStorage) FindFrom(from uint) ([]*model.Post, error) {
	var posts []*model.Post
	if err := us.db.Offset(int(from)).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindFromTo fetches posts from 'from' to 'to' (including 'from' and 'to')
func (us *PostStorage) FindFromTo(from, to uint) ([]*model.Post, error) {
	var posts []*model.Post
	if err := us.db.Offset(int(from)).Limit(int(to - from + 1)).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindByUser finds all posts by user id
func (us *PostStorage) FindByUser(userID uint) ([]*model.Post, error) {
	var posts []*model.Post
	if err := us.db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindByUserFrom fetches posts by user id from 'from'
func (us *PostStorage) FindByUserFrom(userID, from uint) ([]*model.Post, error) {
	var posts []*model.Post
	if err := us.db.Where("user_id = ?", userID).Offset(int(from)).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindByUserFromTo fetches posts by user id from 'from' to 'to' (including 'from' and 'to')
func (us *PostStorage) FindByUserFromTo(userID, from, to uint) ([]*model.Post, error) {
	var posts []*model.Post
	if err := us.db.Where("user_id = ?", userID).Offset(int(from)).Limit(int(to - from + 1)).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindByTag finds all posts by tag id
func (us *PostStorage) FindByTag(tagID uint) ([]*model.Post, error) {
	var posts []*model.Post
	if err := us.db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").Where("post_tags.tag_id = ?", tagID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindByTagFrom fetches posts by tag id from 'from'
func (us *PostStorage) FindByTagFrom(tagID, from uint) ([]*model.Post, error) {
	var posts []*model.Post
	if err := us.db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").Where("post_tags.tag_id = ?", tagID).Offset(int(from)).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindByTagFromTo fetches posts by tag id from 'from' to 'to' (including 'from' and 'to')
func (us *PostStorage) FindByTagFromTo(tagID, from, to uint) ([]*model.Post, error) {
	var posts []*model.Post
	if err := us.db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").Where("post_tags.tag_id = ?", tagID).Offset(int(from)).Limit(int(to - from + 1)).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
