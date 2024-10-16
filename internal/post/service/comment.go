package service

import (
	"github.com/SocialNetworkY/Backend/internal/post/model"
)

type (
	CommentRepo interface {
		Add(comment *model.Comment) error
		Save(comment *model.Comment) error
		Delete(id uint) error
		Find(id uint) (*model.Comment, error)
		FindSome(skip, limit int) ([]*model.Comment, error)
		FindByPost(postID uint, skip, limit int) ([]*model.Comment, error)
		FindByUser(userID uint, skip, limit int) ([]*model.Comment, error)
	}

	CommentService struct {
		repo CommentRepo
	}
)

func NewCommentService(r CommentRepo) *CommentService {
	return &CommentService{r}
}

// Find returns a comment by its ID
func (cs *CommentService) Find(id uint) (*model.Comment, error) {
	return cs.repo.Find(id)
}

// FindByPost returns comments for a post
func (cs *CommentService) FindByPost(postID uint, skip, limit int) ([]*model.Comment, error) {
	return cs.repo.FindByPost(postID, skip, limit)
}

// CommentPost adds a comment to a post
func (cs *CommentService) CommentPost(postID, userID uint, content string) error {
	return cs.repo.Add(&model.Comment{
		UserID:  userID,
		PostID:  postID,
		Content: content,
	})
}

// Edit updates a comment
func (cs *CommentService) Edit(id, userID uint, content string) error {
	comment, err := cs.repo.Find(id)
	if err != nil {
		return err
	}

	comment.Content = content
	comment.Edited = true
	comment.EditedBy = userID
	return cs.repo.Save(comment)
}

// Delete removes a comment
func (cs *CommentService) Delete(id uint) error {
	return cs.repo.Delete(id)
}

func (cs *CommentService) DeleteByPost(postID uint) error {
	comments, err := cs.repo.FindByPost(postID, 0, 0)
	if err != nil {
		return err
	}

	for _, comment := range comments {
		if err := cs.Delete(comment.ID); err != nil {
			return err
		}
	}

	return nil
}

func (cs *CommentService) DeleteByUser(userID uint) error {
	comments, err := cs.repo.FindByUser(userID, 0, 0)
	if err != nil {
		return err
	}

	for _, comment := range comments {
		if err := cs.Delete(comment.ID); err != nil {
			return err
		}
	}

	return nil
}
