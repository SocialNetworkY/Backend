package service

import (
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
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
func (s *CommentService) Find(id uint) (*model.Comment, error) {
	return s.repo.Find(id)
}

// FindByPost returns comments for a post
func (s *CommentService) FindByPost(postID uint, skip, limit int) ([]*model.Comment, error) {
	return s.repo.FindByPost(postID, skip, limit)
}

// CommentPost adds a comment to a post
func (s *CommentService) CommentPost(postID, userID uint, content string) error {
	return s.repo.Add(&model.Comment{
		UserID:  userID,
		PostID:  postID,
		Content: content,
	})
}

// EditComment updates a comment
func (s *CommentService) EditComment(id, userID uint, content string) error {
	comment, err := s.repo.Find(id)
	if err != nil {
		return err
	}

	comment.Content = content
	comment.Edited = true
	comment.EditedBy = userID
	return s.repo.Save(comment)
}

// DeleteComment removes a comment
func (s *CommentService) DeleteComment(id uint) error {
	return s.repo.Delete(id)
}
