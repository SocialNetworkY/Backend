package service

import (
	"log"

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
		Search(query string, skip, limit int) ([]*model.Comment, error)
		Statistic() (*model.CommentStatistic, error)
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
	log.Printf("Finding comment with ID: %d\n", id)
	comment, err := cs.repo.Find(id)
	if err != nil {
		log.Printf("Error finding comment: %v\n", err)
		return nil, err
	}
	log.Printf("Comment found: %v\n", comment)
	return comment, nil
}

// FindByPost returns comments for a post
func (cs *CommentService) FindByPost(postID uint, skip, limit int) ([]*model.Comment, error) {
	log.Printf("Finding comments for post with ID: %d, skip: %d, limit: %d\n", postID, skip, limit)
	comments, err := cs.repo.FindByPost(postID, skip, limit)
	if err != nil {
		log.Printf("Error finding comments: %v\n", err)
		return nil, err
	}
	log.Printf("Comments found: %v\n", comments)
	return comments, nil
}

// CommentPost adds a comment to a post
func (cs *CommentService) CommentPost(postID, userID uint, content string) error {
	log.Printf("Adding comment to post with ID: %d by user with ID: %d\n", postID, userID)
	err := cs.repo.Add(&model.Comment{
		UserID:  userID,
		PostID:  postID,
		Content: content,
	})
	if err != nil {
		log.Printf("Error adding comment: %v\n", err)
		return err
	}
	log.Printf("Comment added successfully to post with ID: %d\n", postID)
	return nil
}

// Edit updates a comment
func (cs *CommentService) Edit(id, userID uint, content string) error {
	log.Printf("Editing comment with ID: %d by user with ID: %d\n", id, userID)
	comment, err := cs.repo.Find(id)
	if err != nil {
		log.Printf("Error finding comment: %v\n", err)
		return err
	}

	comment.Content = content
	comment.Edited = true
	comment.EditedBy = userID
	err = cs.repo.Save(comment)
	if err != nil {
		log.Printf("Error saving comment: %v\n", err)
		return err
	}
	log.Printf("Comment with ID: %d edited successfully\n", id)
	return nil
}

// Delete removes a comment
func (cs *CommentService) Delete(id uint) error {
	log.Printf("Deleting comment with ID: %d\n", id)
	err := cs.repo.Delete(id)
	if err != nil {
		log.Printf("Error deleting comment: %v\n", err)
		return err
	}
	log.Printf("Comment with ID: %d deleted successfully\n", id)
	return nil
}

func (cs *CommentService) DeleteByPost(postID uint) error {
	log.Printf("Deleting comments for post with ID: %d\n", postID)
	comments, err := cs.repo.FindByPost(postID, 0, 0)
	if err != nil {
		log.Printf("Error finding comments: %v\n", err)
		return err
	}

	for _, comment := range comments {
		if err := cs.Delete(comment.ID); err != nil {
			return err
		}
	}
	log.Printf("Comments for post with ID: %d deleted successfully\n", postID)
	return nil
}

func (cs *CommentService) DeleteByUser(userID uint) error {
	log.Printf("Deleting comments by user with ID: %d\n", userID)
	comments, err := cs.repo.FindByUser(userID, 0, 0)
	if err != nil {
		log.Printf("Error finding comments: %v\n", err)
		return err
	}

	for _, comment := range comments {
		if err := cs.Delete(comment.ID); err != nil {
			return err
		}
	}
	log.Printf("Comments by user with ID: %d deleted successfully\n", userID)
	return nil
}

func (cs *CommentService) Search(query string, skip, limit int) ([]*model.Comment, error) {
	log.Printf("Searching comments with query: %s, skip: %d, limit: %d\n", query, skip, limit)
	comments, err := cs.repo.Search(query, skip, limit)
	if err != nil {
		log.Printf("Error searching comments: %v\n", err)
		return nil, err
	}
	log.Printf("Comments found: %v\n", comments)
	return comments, nil
}

func (cs *CommentService) Statistic() (*model.CommentStatistic, error) {
	return cs.repo.Statistic()
}
