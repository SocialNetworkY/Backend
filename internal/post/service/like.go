package service

import (
	"errors"
	"github.com/SocialNetworkY/Backend/internal/post/model"
)

type (
	LikeRepo interface {
		Add(like *model.Like) error
		Delete(id uint) error
		FindByPost(postID uint, skip, limit int) ([]*model.Like, error)
		FindByUser(userID uint, skip, limit int) ([]*model.Like, error)
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
func (ls *LikeService) LikePost(postID, userID uint) error {
	like, _ := ls.repo.FindByPostUser(postID, userID)
	if like != nil {
		return errors.New("user already liked the post")
	}

	return ls.repo.Add(&model.Like{
		UserID: userID,
		PostID: postID,
	})
}

// UnlikePost removes a like from a post
func (ls *LikeService) UnlikePost(postID, userID uint) error {
	like, _ := ls.repo.FindByPostUser(postID, userID)
	if like == nil {
		return errors.New("user has not liked the post")
	}

	return ls.repo.Delete(like.ID)
}

func (ls *LikeService) FindByPostUser(postID, userID uint) (*model.Like, error) {
	return ls.repo.FindByPostUser(postID, userID)
}

func (ls *LikeService) FindByPost(postID uint, skip, limit int) ([]*model.Like, error) {
	return ls.repo.FindByPost(postID, skip, limit)
}

func (ls *LikeService) FindByUser(userID uint, skip, limit int) ([]*model.Like, error) {
	return ls.repo.FindByUser(userID, skip, limit)
}

func (ls *LikeService) Delete(id uint) error {
	return ls.repo.Delete(id)
}

func (ls *LikeService) DeleteByPost(postID uint) error {
	likes, err := ls.repo.FindByPost(postID, 0, 0)
	if err != nil {
		return err
	}

	for _, like := range likes {
		if err := ls.Delete(like.ID); err != nil {
			return err
		}
	}

	return nil
}

func (ls *LikeService) DeleteByUser(userID uint) error {
	likes, err := ls.repo.FindByUser(userID, 0, 0)
	if err != nil {
		return err
	}

	for _, like := range likes {
		if err := ls.Delete(like.ID); err != nil {
			return err
		}
	}

	return nil
}
