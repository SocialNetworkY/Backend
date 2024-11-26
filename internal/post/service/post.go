package service

import (
	"context"
	"log"

	"github.com/SocialNetworkY/Backend/internal/post/model"
)

type (
	PostRepo interface {
		Add(post *model.Post) error
		Save(post *model.Post) error
		Delete(id uint) error
		Find(id uint) (*model.Post, error)
		FindSome(skip, limit int) ([]*model.Post, error)
		FindByUser(userID uint, skip, limit int) ([]*model.Post, error)
		FindByTag(tagID uint, skip, limit int) ([]*model.Post, error)
		ClearAssociations(postID uint) error
		Search(query string, skip, limit int) ([]*model.Post, error)
		Statistic() (*model.PostStatistic, error)
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
	log.Printf("Creating post: %v\n", post)
	if err := ps.processTags(post); err != nil {
		log.Printf("Error processing tags: %v\n", err)
		return err
	}

	if err := ps.repo.Add(post); err != nil {
		log.Printf("Error adding post: %v\n", err)
		return err
	}

	log.Printf("Post created successfully: %v\n", post)
	return nil
}

// Update updates a post
func (ps *PostService) Update(post *model.Post) error {
	log.Printf("Updating post: %v\n", post)
	if err := ps.repo.ClearAssociations(post.ID); err != nil {
		log.Printf("Error clearing associations: %v\n", err)
		return err
	}

	if err := ps.processTags(post); err != nil {
		log.Printf("Error processing tags: %v\n", err)
		return err
	}

	if err := ps.repo.Save(post); err != nil {
		log.Printf("Error saving post: %v\n", err)
		return err
	}

	log.Printf("Post updated successfully: %v\n", post)
	return nil
}

// Find returns a post by its ID
func (ps *PostService) Find(id uint) (*model.Post, error) {
	log.Printf("Finding post with ID: %d\n", id)
	post, err := ps.repo.Find(id)
	if err != nil {
		log.Printf("Error finding post: %v\n", err)
		return nil, err
	}

	log.Printf("Post found: %v\n", post)
	return post, nil
}

// Delete deletes a post by its ID
func (ps *PostService) Delete(id uint) error {
	log.Printf("Deleting post with ID: %d\n", id)
	if err := ps.repo.ClearAssociations(id); err != nil {
		log.Printf("Error clearing associations: %v\n", err)
		return err
	}

	if err := ps.ls.DeleteByPost(id); err != nil {
		log.Printf("Error deleting likes: %v\n", err)
		return err
	}

	if err := ps.cs.DeleteByPost(id); err != nil {
		log.Printf("Error deleting comments: %v\n", err)
		return err
	}

	if err := ps.rg.DeletePostReports(context.Background(), id); err != nil {
		log.Printf("Error deleting reports: %v\n", err)
		return err
	}

	if err := ps.repo.Delete(id); err != nil {
		log.Printf("Error deleting post: %v\n", err)
		return err
	}

	log.Printf("Post deleted successfully: %d\n", id)
	return nil
}

// DeleteByUser deletes all posts by a user
func (ps *PostService) DeleteByUser(userID uint) error {
	log.Printf("Deleting posts by user with ID: %d\n", userID)
	posts, err := ps.repo.FindByUser(userID, 0, 0)
	if err != nil {
		log.Printf("Error finding posts: %v\n", err)
		return err
	}

	for _, post := range posts {
		if err := ps.Delete(post.ID); err != nil {
			return err
		}
	}

	log.Printf("Posts by user with ID %d deleted successfully\n", userID)
	return nil
}

// AddTag adds a tag to a post
func (ps *PostService) AddTag(postID uint, tagName string) error {
	log.Printf("Adding tag %s to post %d\n", tagName, postID)
	post, err := ps.repo.Find(postID)
	if err != nil {
		log.Printf("Error finding post: %v\n", err)
		return err
	}

	tag, err := ps.ts.FindOrCreate(tagName)
	if err != nil {
		log.Printf("Error finding or creating tag: %v\n", err)
		return err
	}

	// Check if the tag is already in the post
	for _, t := range post.Tags {
		if t.ID == tag.ID {
			return nil
		}
	}

	post.Tags = append(post.Tags, tag)

	if err := ps.repo.Save(post); err != nil {
		log.Printf("Error saving post: %v\n", err)
		return err
	}

	log.Printf("Tag %s added to post %d successfully\n", tagName, postID)
	return nil
}

// FindSome returns some posts with pagination parameters
func (ps *PostService) FindSome(skip, limit int) ([]*model.Post, error) {
	log.Printf("Finding some posts with skip: %d and limit: %d\n", skip, limit)
	posts, err := ps.repo.FindSome(skip, limit)
	if err != nil {
		log.Printf("Error finding some posts: %v\n", err)
		return nil, err
	}

	log.Printf("Found some posts: %v\n", posts)
	return posts, nil
}

// FindByUser returns some posts by user ID with pagination parameters
func (ps *PostService) FindByUser(userID uint, skip, limit int) ([]*model.Post, error) {
	log.Printf("Finding posts by user with ID: %d, skip: %d, limit: %d\n", userID, skip, limit)
	posts, err := ps.repo.FindByUser(userID, skip, limit)
	if err != nil {
		log.Printf("Error finding posts by user: %v\n", err)
		return nil, err
	}

	log.Printf("Found posts by user: %v\n", posts)
	return posts, nil
}

// processTags processes tags of a post
func (ps *PostService) processTags(post *model.Post) error {
	log.Printf("Processing tags: %v\n", post.Tags)
	for i, tag := range post.Tags {
		var err error
		if post.Tags[i], err = ps.ts.FindOrCreate(tag.Name); err != nil {
			log.Printf("Error finding or creating tag: %v\n", err)
			return err
		}
	}

	log.Printf("Tags processed successfully: %v\n", post.Tags)
	return nil
}

func (ps *PostService) Search(query string, skip, limit int) ([]*model.Post, error) {
	log.Printf("Searching posts with query: %s, skip: %d, limit: %d\n", query, skip, limit)
	posts, err := ps.repo.Search(query, skip, limit)
	if err != nil {
		log.Printf("Error searching posts: %v\n", err)
		return nil, err
	}

	log.Printf("Found posts: %v\n", posts)
	return posts, nil
}

func (ps *PostService) Statistic() (*model.PostStatistic, error) {
	return ps.repo.Statistic()
}
