package mysql

import (
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
	"gorm.io/gorm"
)

type CommentStorage struct {
	db *gorm.DB
}

func NewCommentStorage(db *gorm.DB) *CommentStorage {
	return &CommentStorage{
		db: db,
	}
}

func (cs *CommentStorage) Add(comment *model.Comment) error {
	if err := cs.db.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func (cs *CommentStorage) Save(comment *model.Comment) error {
	if err := cs.db.Save(comment).Error; err != nil {
		return err
	}
	return nil
}

func (cs *CommentStorage) Delete(id uint) error {
	if err := cs.db.Delete(&model.Comment{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Find finds a comment by id
func (cs *CommentStorage) Find(id uint) (*model.Comment, error) {
	comment := &model.Comment{}
	if err := cs.db.First(comment, id).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// FindAll finds all comments
func (cs *CommentStorage) FindAll() ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := cs.db.Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindFrom fetches comments from 'from'
func (cs *CommentStorage) FindFrom(from uint) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := cs.db.Offset(int(from)).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindFromTo fetches comments from 'from' to 'to'
func (cs *CommentStorage) FindFromTo(from, to uint) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := cs.db.Offset(int(from)).Limit(int(to - from + 1)).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByPost finds all comments by post id
func (cs *CommentStorage) FindByPost(postID uint) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := cs.db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByPostFrom fetches comments by post id from 'from'
func (cs *CommentStorage) FindByPostFrom(postID, from uint) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := cs.db.Where("post_id = ?", postID).Offset(int(from)).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByPostFromTo fetches comments by post id from 'from' to 'to'
func (cs *CommentStorage) FindByPostFromTo(postID, from, to uint) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := cs.db.Where("post_id = ?", postID).Offset(int(from)).Limit(int(to - from + 1)).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByUser finds all comments by user id
func (cs *CommentStorage) FindByUser(userID uint) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := cs.db.Where("user_id = ?", userID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByUserFrom fetches comments by user id from 'from'
func (cs *CommentStorage) FindByUserFrom(userID, from uint) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := cs.db.Where("user_id = ?", userID).Offset(int(from)).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByUserFromTo fetches comments by user id from 'from' to 'to'
func (cs *CommentStorage) FindByUserFromTo(userID, from, to uint) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := cs.db.Where("user_id = ?", userID).Offset(int(from)).Limit(int(to - from + 1)).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
