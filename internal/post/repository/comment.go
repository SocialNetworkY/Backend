package repository

import (
	"log"

	"github.com/SocialNetworkY/Backend/internal/post/model"
	"gorm.io/gorm"
)

type (
	CommentSearch interface {
		Index(comment *model.Comment) error
		Delete(id uint) error
		Search(query string, skip, limit int) ([]uint, error)
	}

	CommentRepository struct {
		db *gorm.DB
		s  CommentSearch
	}
)

func NewCommentRepository(db *gorm.DB, s CommentSearch) *CommentRepository {
	return &CommentRepository{
		db: db,
		s:  s,
	}
}

func (cr *CommentRepository) Add(comment *model.Comment) error {
	log.Printf("Adding comment: %v\n", comment)
	if err := cr.db.Create(comment).Error; err != nil {
		log.Printf("Error adding comment: %v\n", err)
		return err
	}
	if err := cr.s.Index(comment); err != nil {
		log.Printf("Error indexing comment: %v\n", err)
		return err
	}
	log.Printf("Comment added successfully: %v\n", comment)
	return nil
}

func (cr *CommentRepository) Save(comment *model.Comment) error {
	log.Printf("Saving comment: %v\n", comment)
	if err := cr.db.Save(comment).Error; err != nil {
		log.Printf("Error saving comment: %v\n", err)
		return err
	}
	if err := cr.s.Index(comment); err != nil {
		log.Printf("Error indexing comment: %v\n", err)
		return err
	}
	log.Printf("Comment saved successfully: %v\n", comment)
	return nil
}

func (cr *CommentRepository) Delete(id uint) error {
	log.Printf("Deleting comment with ID: %d\n", id)
	if err := cr.db.Delete(&model.Comment{ID: id}).Error; err != nil {
		log.Printf("Error deleting comment: %v\n", err)
		return err
	}
	if err := cr.s.Delete(id); err != nil {
		log.Printf("Error deleting comment from search index: %v\n", err)
		return err
	}
	log.Printf("Comment deleted successfully: %d\n", id)
	return nil
}

func (cr *CommentRepository) Find(id uint) (*model.Comment, error) {
	log.Printf("Finding comment with ID: %d\n", id)
	comment := &model.Comment{}
	if err := cr.db.First(comment, id).Error; err != nil {
		log.Printf("Error finding comment: %v\n", err)
		return nil, err
	}
	log.Printf("Comment found: %v\n", comment)
	return comment, nil
}

func (cr *CommentRepository) FindSome(skip, limit int) ([]*model.Comment, error) {
	log.Printf("Finding some comments with skip: %d and limit: %d\n", skip, limit)
	var comments []*model.Comment
	if limit < 0 {
		skip = -1
	}
	if err := cr.db.Offset(skip).Limit(limit).Find(&comments).Error; err != nil {
		log.Printf("Error finding some comments: %v\n", err)
		return nil, err
	}
	log.Printf("Found some comments: %v\n", comments)
	return comments, nil
}

func (cr *CommentRepository) FindByPost(postID uint, skip, limit int) ([]*model.Comment, error) {
	log.Printf("Finding comments for post with ID: %d, skip: %d, limit: %d\n", postID, skip, limit)
	var comments []*model.Comment
	if limit < 0 {
		skip = -1
	}
	if err := cr.db.Offset(skip).Where("post_id = ?", postID).Limit(limit).Find(&comments).Error; err != nil {
		log.Printf("Error finding comments: %v\n", err)
		return nil, err
	}
	log.Printf("Comments found: %v\n", comments)
	return comments, nil
}

func (cr *CommentRepository) FindByUser(userID uint, skip, limit int) ([]*model.Comment, error) {
	log.Printf("Finding comments by user with ID: %d, skip: %d, limit: %d\n", userID, skip, limit)
	var comments []*model.Comment
	if limit < 0 {
		skip = -1
	}
	if err := cr.db.Offset(skip).Where("user_id = ?", userID).Limit(limit).Find(&comments).Error; err != nil {
		log.Printf("Error finding comments: %v\n", err)
		return nil, err
	}
	log.Printf("Comments found: %v\n", comments)
	return comments, nil
}

func (cr *CommentRepository) Search(query string, skip, limit int) ([]*model.Comment, error) {
	log.Printf("Searching comments with query: %s, skip: %d, limit: %d\n", query, skip, limit)
	ids, err := cr.s.Search(query, skip, limit)
	if err != nil {
		log.Printf("Error searching comments: %v\n", err)
		return nil, err
	}
	var comments []*model.Comment
	if err := cr.db.Where("id IN ?", ids).Find(&comments).Error; err != nil {
		log.Printf("Error finding comments by IDs: %v\n", err)
		return nil, err
	}
	log.Printf("Comments found: %v\n", comments)
	return comments, nil
}

func (cr *CommentRepository) Statistic() (*model.CommentStatistic, error) {
	log.Println("Getting comment statistics")

	var stat model.CommentStatistic
	err := cr.db.Model(&model.Comment{}).
		Select("COUNT(*) AS total, " +
			"SUM(CASE WHEN edited_by != 0 THEN 1 ELSE 0 END) AS edited").
		Scan(&stat).Error

	if err != nil {
		log.Printf("Error getting comment statistics: %v\n", err)
		return nil, err
	}

	log.Printf("Comment statistics found: %+v\n", stat)
	return &stat, nil
}
