package service

import (
	"errors"
	"log"

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
	log.Printf("Adding like to post with ID: %d by user with ID: %d\n", postID, userID)
	like, _ := ls.repo.FindByPostUser(postID, userID)
	if like != nil {
		log.Printf("Error: user already liked the post\n")
		return errors.New("user already liked the post")
	}

	err := ls.repo.Add(&model.Like{
		UserID: userID,
		PostID: postID,
	})
	if err != nil {
		log.Printf("Error adding like: %v\n", err)
		return err
	}

	log.Printf("Like added successfully to post with ID: %d by user with ID: %d\n", postID, userID)
	return nil
}

// UnlikePost removes a like from a post
func (ls *LikeService) UnlikePost(postID, userID uint) error {
	log.Printf("Removing like from post with ID: %d by user with ID: %d\n", postID, userID)
	like, _ := ls.repo.FindByPostUser(postID, userID)
	if like == nil {
		log.Printf("Error: user has not liked the post\n")
		return errors.New("user has not liked the post")
	}

	err := ls.repo.Delete(like.ID)
	if err != nil {
		log.Printf("Error removing like: %v\n", err)
		return err
	}

	log.Printf("Like removed successfully from post with ID: %d by user with ID: %d\n", postID, userID)
	return nil
}

func (ls *LikeService) FindByPostUser(postID, userID uint) (*model.Like, error) {
	log.Printf("Finding like for post with ID: %d by user with ID: %d\n", postID, userID)
	like, err := ls.repo.FindByPostUser(postID, userID)
	if err != nil {
		log.Printf("Error finding like: %v\n", err)
		return nil, err
	}
	log.Printf("Like found: %v\n", like)
	return like, nil
}

func (ls *LikeService) FindByPost(postID uint, skip, limit int) ([]*model.Like, error) {
	log.Printf("Finding likes for post with ID: %d, skip: %d, limit: %d\n", postID, skip, limit)
	likes, err := ls.repo.FindByPost(postID, skip, limit)
	if err != nil {
		log.Printf("Error finding likes: %v\n", err)
		return nil, err
	}
	log.Printf("Likes found: %v\n", likes)
	return likes, nil
}

func (ls *LikeService) FindByUser(userID uint, skip, limit int) ([]*model.Like, error) {
	log.Printf("Finding likes by user with ID: %d, skip: %d, limit: %d\n", userID, skip, limit)
	likes, err := ls.repo.FindByUser(userID, skip, limit)
	if err != nil {
		log.Printf("Error finding likes: %v\n", err)
		return nil, err
	}
	log.Printf("Likes found: %v\n", likes)
	return likes, nil
}

func (ls *LikeService) Delete(id uint) error {
	log.Printf("Deleting like with ID: %d\n", id)
	err := ls.repo.Delete(id)
	if err != nil {
		log.Printf("Error deleting like: %v\n", err)
		return err
	}
	log.Printf("Like with ID: %d deleted successfully\n", id)
	return nil
}

func (ls *LikeService) DeleteByPost(postID uint) error {
	log.Printf("Deleting likes for post with ID: %d\n", postID)
	likes, err := ls.repo.FindByPost(postID, 0, 0)
	if err != nil {
		log.Printf("Error finding likes: %v\n", err)
		return err
	}

	for _, like := range likes {
		if err := ls.Delete(like.ID); err != nil {
			return err
		}
	}

	log.Printf("Likes for post with ID: %d deleted successfully\n", postID)
	return nil
}

func (ls *LikeService) DeleteByUser(userID uint) error {
	log.Printf("Deleting likes by user with ID: %d\n", userID)
	likes, err := ls.repo.FindByUser(userID, 0, 0)
	if err != nil {
		log.Printf("Error finding likes: %v\n", err)
		return err
	}

	for _, like := range likes {
		if err := ls.Delete(like.ID); err != nil {
			return err
		}
	}

	log.Printf("Likes by user with ID: %d deleted successfully\n", userID)
	return nil
}
