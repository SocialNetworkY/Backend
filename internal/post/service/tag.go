package service

import (
	"github.com/SocialNetworkY/Backend/internal/post/model"
)

type (
	TagRepo interface {
		Add(tag *model.Tag) error
		Save(tag *model.Tag) error
		Delete(id uint) error
		Find(id uint) (*model.Tag, error)
		FindByName(name string) (*model.Tag, error)
		FindSome(skip, limit int) ([]*model.Tag, error)
		FindByPost(postID uint, skip, limit int) ([]*model.Tag, error)
		ClearAssociations(tagID uint) error
	}

	TagService struct {
		repo TagRepo
	}
)

func NewTagService(r TagRepo) *TagService {
	return &TagService{
		repo: r,
	}
}

// Add adds a new tag
func (ts *TagService) Add(tag *model.Tag) error {
	return ts.repo.Add(tag)
}

// Delete deletes a tag by id
func (ts *TagService) Delete(id uint) error {
	if err := ts.repo.ClearAssociations(id); err != nil {
		return err
	}

	return ts.repo.Delete(id)
}

// DeleteByName deletes a tag by name
func (ts *TagService) DeleteByName(name string) error {
	tag, err := ts.repo.FindByName(name)
	if err != nil {
		return err
	}

	return ts.Delete(tag.ID)
}

// Find finds a tag by id
func (ts *TagService) Find(id uint) (*model.Tag, error) {
	return ts.repo.Find(id)
}

// FindByName finds a tag by name
func (ts *TagService) FindByName(name string) (*model.Tag, error) {
	return ts.repo.FindByName(name)
}

// Exists checks if a tag exists
func (ts *TagService) Exists(id uint) (bool, error) {
	_, err := ts.repo.Find(id)
	return err == nil, err
}

// ExistsByName checks if a tag exists
func (ts *TagService) ExistsByName(name string) bool {
	_, err := ts.repo.FindByName(name)
	return err == nil
}

// FindOrCreate finds a tag by name or creates a new one
func (ts *TagService) FindOrCreate(name string) (*model.Tag, error) {
	tag, err := ts.repo.FindByName(name)
	if err != nil {
		tag = &model.Tag{
			Name: name,
		}

		if err := ts.repo.Add(tag); err != nil {
			return nil, err
		}
	}

	return tag, nil
}
