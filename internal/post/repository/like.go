package repository

import (
	"log"

	"github.com/SocialNetworkY/Backend/internal/post/model"
	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{
		db: db,
	}
}

func (lr *LikeRepository) Add(like *model.Like) error {
	log.Printf("Adding like: %v\n", like)
	if err := lr.db.Create(like).Error; err != nil {
		log.Printf("Error adding like: %v\n", err)
		return err
	}
	log.Printf("Like added successfully: %v\n", like)
	return nil
}

func (lr *LikeRepository) Save(like *model.Like) error {
	log.Printf("Saving like: %v\n", like)
	if err := lr.db.Save(like).Error; err != nil {
		log.Printf("Error saving like: %v\n", err)
		return err
	}
	log.Printf("Like saved successfully: %v\n", like)
	return nil
}

func (lr *LikeRepository) Delete(id uint) error {
	log.Printf("Deleting like with ID: %d\n", id)
	if err := lr.db.Delete(&model.Like{ID: id}).Error; err != nil {
		log.Printf("Error deleting like: %v\n", err)
		return err
	}
	log.Printf("Like deleted successfully: %d\n", id)
	return nil
}

func (lr *LikeRepository) Find(id uint) (*model.Like, error) {
	log.Printf("Finding like with ID: %d\n", id)
	like := &model.Like{}
	if err := lr.db.First(like, id).Error; err != nil {
		log.Printf("Error finding like: %v\n", err)
		return nil, err
	}
	log.Printf("Like found: %v\n", like)
	return like, nil
}

func (lr *LikeRepository) FindSome(skip, limit int) ([]*model.Like, error) {
	log.Printf("Finding some likes with skip: %d and limit: %d\n", skip, limit)
	var likes []*model.Like
	if limit < 0 {
		skip = -1
	}
	if err := lr.db.Offset(skip).Limit(limit).Find(&likes).Error; err != nil {
		log.Printf("Error finding some likes: %v\n", err)
		return nil, err
	}
	log.Printf("Found some likes: %v\n", likes)
	return likes, nil
}

func (lr *LikeRepository) FindByPost(postID uint, skip, limit int) ([]*model.Like, error) {
	log.Printf("Finding likes for post with ID: %d, skip: %d, limit: %d\n", postID, skip, limit)
	var likes []*model.Like
	if limit < 0 {
		skip = -1
	}
	if err := lr.db.Offset(skip).Where("post_id = ?", postID).Limit(limit).Find(&likes).Error; err != nil {
		log.Printf("Error finding likes: %v\n", err)
		return nil, err
	}
	log.Printf("Likes found: %v\n", likes)
	return likes, nil
}

func (lr *LikeRepository) FindByUser(userID uint, skip, limit int) ([]*model.Like, error) {
	log.Printf("Finding likes by user with ID: %d, skip: %d, limit: %d\n", userID, skip, limit)
	var likes []*model.Like
	if limit < 0 {
		skip = -1
	}
	if err := lr.db.Offset(skip).Where("user_id = ?", userID).Limit(limit).Find(&likes).Error; err != nil {
		log.Printf("Error finding likes: %v\n", err)
		return nil, err
	}
	log.Printf("Likes found: %v\n", likes)
	return likes, nil
}

func (lr *LikeRepository) FindByPostUser(postID, userID uint) (*model.Like, error) {
	log.Printf("Finding like for post with ID: %d by user with ID: %d\n", postID, userID)
	var like *model.Like
	if err := lr.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error; err != nil {
		log.Printf("Error finding like: %v\n", err)
		return nil, err
	}
	log.Printf("Like found: %v\n", like)
	return like, nil
}
