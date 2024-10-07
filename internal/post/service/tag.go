package service

import (
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
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
func (ts *TagService) Add(tag *model.Tag) error {
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

// Find finds a tag by id
func (ts *TagService) Find(id uint) (*model.Tag, error) {
	return ts.s.Find(id)
}

// FindByName finds a tag by name
func (ts *TagService) FindByName(name string) (*model.Tag, error) {
	return ts.s.FindByName(name)
}

// Exists checks if a tag exists
func (ts *TagService) Exists(id uint) (bool, error) {
	_, err := ts.s.Find(id)
	return err == nil, err
}

// ExistsByName checks if a tag exists
func (ts *TagService) ExistsByName(name string) bool {
	_, err := ts.s.FindByName(name)
	return err == nil
}

// FindOrCreate finds a tag by name or creates a new one
func (ts *TagService) FindOrCreate(name string) (*model.Tag, error) {
	tag, err := ts.s.FindByName(name)
	if err != nil {
		tag = &model.Tag{
			Name: name,
		}

		if err := ts.s.Add(tag); err != nil {
			return nil, err
		}
	}

	return tag, nil
}
