package service

import (
	"errors"
	"github.com/SocialNetworkY/Backend/internal/post/model"
)

type (
	LikeRepo interface {
		Add(like *model.Like) error
		Delete(id uint) error
		FindByPostUser(postID, userID uint) (*model.Like, error)
	}

	LikeService struct {
		repo LikeRepo
	}
)

func NewLikeService(r LikeRepo) *LikeService {
	return &LikeService{r}
}

// LikePost adds a like to a post
func (s *LikeService) LikePost(postID, userID uint) error {
	like, _ := s.repo.FindByPostUser(postID, userID)
	if like != nil {
		return errors.New("user already liked the post")
	}

	return s.repo.Add(&model.Like{
		UserID: userID,
		PostID: postID,
	})
}

// UnlikePost removes a like from a post
func (s *LikeService) UnlikePost(postID, userID uint) error {
	like, _ := s.repo.FindByPostUser(postID, userID)
	if like == nil {
		return errors.New("user has not liked the post")
	}

	return s.repo.Delete(like.ID)
}
