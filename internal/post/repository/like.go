package repository

import (
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{
		db: db,
	}
}

func (lr *LikeRepository) Add(like *model.Like) error {
	if err := lr.db.Create(like).Error; err != nil {
		return err
	}
	return nil
}

func (lr *LikeRepository) Save(like *model.Like) error {
	if err := lr.db.Save(like).Error; err != nil {
		return err
	}
	return nil
}

func (lr *LikeRepository) Delete(id uint) error {
	if err := lr.db.Delete(&model.Like{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

// Find finds a like by id
func (lr *LikeRepository) Find(id uint) (*model.Like, error) {
	like := &model.Like{}
	if err := lr.db.First(like, id).Error; err != nil {
		return nil, err
	}
	return like, nil
}

// FindSome fetches some likes
func (lr *LikeRepository) FindSome(skip, limit int) ([]*model.Like, error) {
	var likes []*model.Like
	if err := lr.db.Offset(skip).Limit(limit).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

// FindByPost finds some likes by post id
func (lr *LikeRepository) FindByPost(postID uint, skip, limit int) ([]*model.Like, error) {
	var likes []*model.Like
	if err := lr.db.Offset(skip).Limit(limit).Where("post_id = ?", postID).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

// FindByUser finds some likes by user id
func (lr *LikeRepository) FindByUser(userID uint, skip, limit int) ([]*model.Like, error) {
	var likes []*model.Like
	if err := lr.db.Offset(skip).Limit(limit).Where("user_id = ?", userID).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

// FindByPostUser finds like by post id and user id
func (lr *LikeRepository) FindByPostUser(postID, userID uint) (*model.Like, error) {
	var like *model.Like
	if err := lr.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error; err != nil {
		return nil, err
	}
	return like, nil
}
