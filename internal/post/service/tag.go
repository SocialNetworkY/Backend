package service

import (
	"log"

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
		Search(query string, skip, limit int) ([]*model.Tag, error)
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
	log.Printf("Adding tag: %v\n", tag)
	err := ts.repo.Add(tag)
	if err != nil {
		log.Printf("Error adding tag: %v\n", err)
		return err
	}
	log.Printf("Tag added successfully: %v\n", tag)
	return nil
}

// Delete deletes a tag by id
func (ts *TagService) Delete(id uint) error {
	log.Printf("Deleting tag with ID: %d\n", id)
	if err := ts.repo.ClearAssociations(id); err != nil {
		log.Printf("Error clearing associations: %v\n", err)
		return err
	}

	if err := ts.repo.Delete(id); err != nil {
		log.Printf("Error deleting tag: %v\n", err)
		return err
	}

	log.Printf("Tag deleted successfully: %d\n", id)
	return nil
}

// DeleteByName deletes a tag by name
func (ts *TagService) DeleteByName(name string) error {
	log.Printf("Deleting tag with name: %s\n", name)
	tag, err := ts.repo.FindByName(name)
	if err != nil {
		log.Printf("Error finding tag: %v\n", err)
		return err
	}

	if err := ts.Delete(tag.ID); err != nil {
		return err
	}

	log.Printf("Tag with name %s deleted successfully\n", name)
	return nil
}

// Find finds a tag by id
func (ts *TagService) Find(id uint) (*model.Tag, error) {
	log.Printf("Finding tag with ID: %d\n", id)
	tag, err := ts.repo.Find(id)
	if err != nil {
		log.Printf("Error finding tag: %v\n", err)
		return nil, err
	}
	log.Printf("Tag found: %v\n", tag)
	return tag, nil
}

// FindByName finds a tag by name
func (ts *TagService) FindByName(name string) (*model.Tag, error) {
	log.Printf("Finding tag with name: %s\n", name)
	tag, err := ts.repo.FindByName(name)
	if err != nil {
		log.Printf("Error finding tag: %v\n", err)
		return nil, err
	}
	log.Printf("Tag found: %v\n", tag)
	return tag, nil
}

// Exists checks if a tag exists
func (ts *TagService) Exists(id uint) (bool, error) {
	log.Printf("Checking if tag with ID: %d exists\n", id)
	_, err := ts.repo.Find(id)
	exists := err == nil
	if !exists {
		log.Printf("Tag with ID: %d does not exist\n", id)
	}
	return exists, err
}

// ExistsByName checks if a tag exists
func (ts *TagService) ExistsByName(name string) bool {
	log.Printf("Checking if tag with name: %s exists\n", name)
	_, err := ts.repo.FindByName(name)
	exists := err == nil
	if !exists {
		log.Printf("Tag with name: %s does not exist\n", name)
	}
	return exists
}

// FindOrCreate finds a tag by name or creates a new one
func (ts *TagService) FindOrCreate(name string) (*model.Tag, error) {
	log.Printf("Finding or creating tag: %s\n", name)
	tag, err := ts.repo.FindByName(name)
	if err != nil {
		tag = &model.Tag{
			Name: name,
		}

		if err := ts.repo.Add(tag); err != nil {
			log.Printf("Error adding tag: %v\n", err)
			return nil, err
		}
	}

	log.Printf("Tag found or created: %v\n", tag)
	return tag, nil
}

func (ts *TagService) Search(query string, skip, limit int) ([]*model.Tag, error) {
	log.Printf("Searching tags with query: %s, skip: %d, limit: %d\n", query, skip, limit)
	tags, err := ts.repo.Search(query, skip, limit)
	if err != nil {
		log.Printf("Error searching tags: %v\n", err)
		return nil, err
	}
	log.Printf("Tags found: %v\n", tags)
	return tags, nil
}
