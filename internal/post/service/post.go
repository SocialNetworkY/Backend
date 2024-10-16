package service

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/post/model"
)

type (
	PostRepo interface {
		Add(post *model.Post) error
		Save(post *model.Post) error
		Delete(id uint) error
		DeleteByUser(userID uint) error
		Find(id uint) (*model.Post, error)
		FindSome(skip, limit int) ([]*model.Post, error)
		FindByUser(userID uint, skip, limit int) ([]*model.Post, error)
		FindByTag(tagID uint, skip, limit int) ([]*model.Post, error)
		ClearAssociations(postID uint) error
	}

	ReportGateway interface {
		DeletePostReports(ctx context.Context, postID uint) error
	}

	PostService struct {
		repo PostRepo
		rg   ReportGateway
		ts   *TagService
		cs   *CommentService
		ls   *LikeService
	}
)

func NewPostService(r PostRepo, rg ReportGateway, ts *TagService, cs *CommentService, ls *LikeService) *PostService {
	return &PostService{
		repo: r,
		rg:   rg,
		ts:   ts,
		cs:   cs,
		ls:   ls,
	}
}

// Create creates a new post
func (ps *PostService) Create(post *model.Post) error {
	if err := ps.processTags(post); err != nil {
		return err
	}

	return ps.repo.Add(post)
}

// Update updates a post
func (ps *PostService) Update(post *model.Post) error {
	if err := ps.processTags(post); err != nil {
		return err
	}

	return ps.repo.Save(post)
}

// Find returns a post by its ID
func (ps *PostService) Find(id uint) (*model.Post, error) {
	return ps.repo.Find(id)
}

// Delete deletes a post by its ID
func (ps *PostService) Delete(id uint) error {
	if err := ps.repo.ClearAssociations(id); err != nil {
		return err
	}

	if err := ps.ls.DeleteByPost(id); err != nil {
		return err
	}

	if err := ps.cs.DeleteByPost(id); err != nil {
		return err
	}

	if err := ps.rg.DeletePostReports(context.Background(), id); err != nil {
		return err
	}

	return ps.repo.Delete(id)
}

// DeleteByUser deletes all posts by a user
func (ps *PostService) DeleteByUser(userID uint) error {
	posts, err := ps.repo.FindByUser(userID, 0, 0)
	if err != nil {
		return err
	}

	for _, post := range posts {
		if err := ps.Delete(post.ID); err != nil {
			return err
		}
	}

	return nil
}

// AddTag adds a tag to a post
func (ps *PostService) AddTag(postID uint, tagName string) error {
	post, err := ps.repo.Find(postID)
	if err != nil {
		return err
	}

	tag, err := ps.ts.FindOrCreate(tagName)
	if err != nil {
		return err
	}

	post.Tags = append(post.Tags, tag)

	return ps.repo.Save(post)
}

// FindSome returns some posts with pagination parameters
func (ps *PostService) FindSome(skip, limit int) ([]*model.Post, error) {
	return ps.repo.FindSome(skip, limit)
}

// FindByUser returns some posts by user ID with pagination parameters
func (ps *PostService) FindByUser(userID uint, skip, limit int) ([]*model.Post, error) {
	return ps.repo.FindByUser(userID, skip, limit)
}

// processTags processes tags of a post
func (ps *PostService) processTags(post *model.Post) error {
	if err := ps.repo.ClearAssociations(post.ID); err != nil {
		return err
	}

	for i, tag := range post.Tags {
		var err error
		if post.Tags[i], err = ps.ts.FindOrCreate(tag.Name); err != nil {
			return err
		}
	}

	return nil
}
