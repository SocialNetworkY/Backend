package service

import (
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
)

type (
	PostStorage interface {
		Add(post *model.Post) error
		Save(post *model.Post) error
		Delete(id uint) error
		Find(id uint) (*model.Post, error)
		FindAll() ([]*model.Post, error)
		FindFrom(from uint) ([]*model.Post, error)
		FindFromTo(from, to uint) ([]*model.Post, error)
		FindByUser(userID uint) ([]*model.Post, error)
		FindByUserFrom(userID, from uint) ([]*model.Post, error)
		FindByUserFromTo(userID, from, to uint) ([]*model.Post, error)
		FindByTag(tagID uint) ([]*model.Post, error)
		FindByTagFrom(tagID, from uint) ([]*model.Post, error)
		FindByTagFromTo(tagID, from, to uint) ([]*model.Post, error)
	}

	PostService struct {
		s  PostStorage
		ts *TagService
	}
)

func NewPostService(s PostStorage, ts *TagService) *PostService {
	return &PostService{
		s:  s,
		ts: ts,
	}
}

// Create creates a new post
func (ps *PostService) Create(post *model.Post) error {
	if err := ps.processTags(post); err != nil {
		return err
	}

	return ps.s.Add(post)
}

// Update updates a post
func (ps *PostService) Update(post *model.Post) error {
	if err := ps.processTags(post); err != nil {
		return err
	}

	return ps.s.Save(post)
}

// Find returns a post by its ID
func (ps *PostService) Find(id uint) (*model.Post, error) {
	return ps.s.Find(id)
}

// Delete deletes a post by its ID
func (ps *PostService) Delete(id uint) error {
	return ps.s.Delete(id)
}

// AddTag adds a tag to a post
func (ps *PostService) AddTag(postID uint, tagName string) error {
	post, err := ps.s.Find(postID)
	if err != nil {
		return err
	}

	tag, err := ps.ts.FindOrCreate(tagName)
	if err != nil {
		return err
	}

	post.Tags = append(post.Tags, tag)

	return ps.s.Save(post)
}

// processTags processes tags of a post
func (ps *PostService) processTags(post *model.Post) error {
	for _, tag := range post.Tags {
		var err error
		if tag, err = ps.ts.FindOrCreate(tag.Name); err != nil {
			return err
		}
	}

	return nil
}
