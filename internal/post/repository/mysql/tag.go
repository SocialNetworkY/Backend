package mysql

import (
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
	"gorm.io/gorm"
)

type TagStorage struct {
	db *gorm.DB
}

func NewTagStorage(db *gorm.DB) *TagStorage {
	return &TagStorage{
		db: db,
	}
}

func (ts *TagStorage) Add(tag *model.Tag) error {
	if err := ts.db.Create(tag).Error; err != nil {
		return err
	}
	return nil
}

func (ts *TagStorage) Save(tag *model.Tag) error {
	if err := ts.db.Save(tag).Error; err != nil {
		return err
	}
	return nil
}

func (ts *TagStorage) Delete(id uint) error {
	if err := ts.db.Delete(&model.Tag{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Find finds a tag by id
func (ts *TagStorage) Find(id uint) (*model.Tag, error) {
	tag := &model.Tag{}
	if err := ts.db.First(tag, id).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

// FindByName finds a tag by name
func (ts *TagStorage) FindByName(name string) (*model.Tag, error) {
	tag := &model.Tag{}
	if err := ts.db.Where("name = ?", name).First(tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

// FindAll finds all tags
func (ts *TagStorage) FindAll() ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := ts.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// FindFrom fetches tags from 'from'
func (ts *TagStorage) FindFrom(from uint) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := ts.db.Offset(int(from)).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// FindFromTo fetches tags from 'from' to 'to'
func (ts *TagStorage) FindFromTo(from, to uint) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := ts.db.Offset(int(from)).Limit(int(to - from + 1)).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// FindByPost finds all tags by post id
func (ts *TagStorage) FindByPost(postID uint) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := ts.db.Joins("JOIN post_tags ON post_tags.tag_id = tags.id").Where("post_tags.post_id = ?", postID).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// FindByPostFrom fetches tags by post id from 'from'
func (ts *TagStorage) FindByPostFrom(postID, from uint) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := ts.db.Joins("JOIN post_tags ON post_tags.tag_id = tags.id").Where("post_tags.post_id = ?", postID).Offset(int(from)).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// FindByPostFromTo fetches tags by post id from 'from' to 'to'
func (ts *TagStorage) FindByPostFromTo(postID, from, to uint) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := ts.db.Joins("JOIN post_tags ON post_tags.tag_id = tags.id").Where("post_tags.post_id = ?", postID).Offset(int(from)).Limit(int(to - from + 1)).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
