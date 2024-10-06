package mysql

import (
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
	"gorm.io/gorm"
)

type LikeStorage struct {
	db *gorm.DB
}

func NewLikeStorage(db *gorm.DB) *LikeStorage {
	return &LikeStorage{
		db: db,
	}
}

func (ls *LikeStorage) Add(like *model.Like) error {
	if err := ls.db.Create(like).Error; err != nil {
		return err
	}
	return nil
}

func (ls *LikeStorage) Save(like *model.Like) error {
	if err := ls.db.Save(like).Error; err != nil {
		return err
	}
	return nil
}

func (ls *LikeStorage) Delete(id uint) error {
	if err := ls.db.Delete(&model.Like{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Find finds a like by id
func (ls *LikeStorage) Find(id uint) (*model.Like, error) {
	like := &model.Like{}
	if err := ls.db.First(like, id).Error; err != nil {
		return nil, err
	}
	return like, nil
}

// FindAll finds all likes
func (ls *LikeStorage) FindAll() ([]*model.Like, error) {
	var likes []*model.Like
	if err := ls.db.Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

// FindFrom fetches likes from 'from'
func (ls *LikeStorage) FindFrom(from uint) ([]*model.Like, error) {
	var likes []*model.Like
	if err := ls.db.Offset(int(from)).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

// FindFromTo fetches likes from 'from' to 'to'
func (ls *LikeStorage) FindFromTo(from, to uint) ([]*model.Like, error) {
	var likes []*model.Like
	if err := ls.db.Offset(int(from)).Limit(int(to - from + 1)).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

// FindByPost finds all likes by post id
func (ls *LikeStorage) FindByPost(postID uint) ([]*model.Like, error) {
	var likes []*model.Like
	if err := ls.db.Where("post_id = ?", postID).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

// FindByPostFrom fetches likes by post id from 'from'
func (ls *LikeStorage) FindByPostFrom(postID, from uint) ([]*model.Like, error) {
	var likes []*model.Like
	if err := ls.db.Where("post_id = ?", postID).Offset(int(from)).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

// FindByPostFromTo fetches likes by post id from 'from' to 'to'
func (ls *LikeStorage) FindByPostFromTo(postID, from, to uint) ([]*model.Like, error) {
	var likes []*model.Like
	if err := ls.db.Where("post_id = ?", postID).Offset(int(from)).Limit(int(to - from + 1)).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

// FindByUser finds all likes by user id
func (ls *LikeStorage) FindByUser(userID uint) ([]*model.Like, error) {
	var likes []*model.Like
	if err := ls.db.Where("user_id = ?", userID).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

// FindByUserFrom fetches likes by user id from 'from'
func (ls *LikeStorage) FindByUserFrom(userID, from uint) ([]*model.Like, error) {
	var likes []*model.Like
	if err := ls.db.Where("user_id = ?", userID).Offset(int(from)).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

// FindByUserFromTo fetches likes by user id from 'from' to 'to'
func (ls *LikeStorage) FindByUserFromTo(userID, from, to uint) ([]*model.Like, error) {
	var likes []*model.Like
	if err := ls.db.Where("user_id = ?", userID).Offset(int(from)).Limit(int(to - from + 1)).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}
