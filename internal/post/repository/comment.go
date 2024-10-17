package repository

import (
	"github.com/SocialNetworkY/Backend/internal/post/model"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (cr *CommentRepository) Add(comment *model.Comment) error {
	if err := cr.db.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func (cr *CommentRepository) Save(comment *model.Comment) error {
	if err := cr.db.Save(comment).Error; err != nil {
		return err
	}
	return nil
}

func (cr *CommentRepository) Delete(id uint) error {
	if err := cr.db.Delete(&model.Comment{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

// Find finds a comment by id
func (cr *CommentRepository) Find(id uint) (*model.Comment, error) {
	comment := &model.Comment{}
	if err := cr.db.First(comment, id).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// FindSome fetches some comments
func (cr *CommentRepository) FindSome(skip, limit int) ([]*model.Comment, error) {
	var comments []*model.Comment
	if limit < 0 {
		skip = -1
	}
	if err := cr.db.Offset(skip).Limit(limit).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByPost finds some comments by post id
func (cr *CommentRepository) FindByPost(postID uint, skip, limit int) ([]*model.Comment, error) {
	var comments []*model.Comment
	if limit < 0 {
		skip = -1
	}
	if err := cr.db.Offset(skip).Where("post_id = ?", postID).Limit(limit).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByUser finds some comments by user id
func (cr *CommentRepository) FindByUser(userID uint, skip, limit int) ([]*model.Comment, error) {
	var comments []*model.Comment
	if limit < 0 {
		skip = -1
	}
	if err := cr.db.Offset(skip).Where("user_id = ?", userID).Limit(limit).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
