package service

import (
	"errors"
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
	"gorm.io/gorm"
)

type (
	TagStorage interface {
		Add(tag *model.Tag) error
		Save(tag *model.Tag) error
		Delete(id uint) error
		Find(id uint) (*model.Tag, error)
		FindByName(name string) (*model.Tag, error)
		FindAll() ([]*model.Tag, error)
		FindFrom(from uint) ([]*model.Tag, error)
		FindFromTo(from, to uint) ([]*model.Tag, error)
		FindByPost(postID uint) ([]*model.Tag, error)
		FindByPostFrom(postID, from uint) ([]*model.Tag, error)
		FindByPostFromTo(postID, from, to uint) ([]*model.Tag, error)
	}

	TagService struct {
		s TagStorage
	}
)

func NewTagService(s TagStorage) *TagService {
	return &TagService{
		s: s,
	}
}

// Add adds a new tag
func (ts *TagService) Add(name string) error {
	// Check if tag exists
	tag, err := ts.s.FindByName(name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return ts.s.Add(tag)
}

// Delete deletes a tag by id
func (ts *TagService) Delete(id uint) error {
	return ts.s.Delete(id)
}

// DeleteByName deletes a tag by name
func (ts *TagService) DeleteByName(name string) error {
	tag, err := ts.s.FindByName(name)
	if err != nil {
		return err
	}

	return ts.s.Delete(tag.ID)
}
