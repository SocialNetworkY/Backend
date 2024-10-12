package repository

import (
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		db: db,
	}
}

func (tr *TagRepository) Add(tag *model.Tag) error {
	if err := tr.db.Create(tag).Error; err != nil {
		return err
	}
	return nil
}

func (tr *TagRepository) Save(tag *model.Tag) error {
	if err := tr.db.Save(tag).Error; err != nil {
		return err
	}
	return nil
}

func (tr *TagRepository) Delete(id uint) error {
	if err := tr.db.Delete(&model.Tag{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

// Find finds a tag by id
func (tr *TagRepository) Find(id uint) (*model.Tag, error) {
	tag := &model.Tag{}
	if err := tr.db.Preload("Posts").First(tag, id).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

// FindByName finds a tag by name
func (tr *TagRepository) FindByName(name string) (*model.Tag, error) {
	tag := &model.Tag{}
	if err := tr.db.Preload("Posts").Where("name = ?", name).First(tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

// FindSome fetches some tags
func (tr *TagRepository) FindSome(skip, limit int) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := tr.db.Preload("Posts").Offset(skip).Limit(limit).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// FindByPost finds some tags by post id
func (tr *TagRepository) FindByPost(postID uint, skip, limit int) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := tr.db.Preload("Posts").Offset(skip).Limit(limit).Joins("JOIN post_tags ON post_tags.tag_id = tags.id").Where("post_tags.post_id = ?", postID).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
