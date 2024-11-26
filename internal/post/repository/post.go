package repository

import (
	"log"
	"sort"

	"github.com/SocialNetworkY/Backend/internal/post/model"
	"gorm.io/gorm"
)

type (
	PostSearch interface {
		Index(post *model.Post) error
		Delete(id uint) error
		Search(query string, skip, limit int) ([]uint, error)
	}

	PostRepository struct {
		db *gorm.DB
		s  PostSearch
	}
)

func NewPostRepository(db *gorm.DB, s PostSearch) *PostRepository {
	return &PostRepository{
		db: db,
		s:  s,
	}
}

func (pr *PostRepository) Add(post *model.Post) error {
	log.Printf("Adding post: %v\n", post)
	if err := pr.db.Create(post).Error; err != nil {
		log.Printf("Error adding post: %v\n", err)
		return err
	}
	if err := pr.s.Index(post); err != nil {
		log.Printf("Error indexing post: %v\n", err)
		return err
	}
	log.Printf("Post added successfully: %v\n", post)
	return nil
}

func (pr *PostRepository) Save(post *model.Post) error {
	log.Printf("Saving post: %v\n", post)
	if err := pr.db.Save(post).Error; err != nil {
		log.Printf("Error saving post: %v\n", err)
		return err
	}
	if err := pr.s.Index(post); err != nil {
		log.Printf("Error indexing post: %v\n", err)
		return err
	}
	log.Printf("Post saved successfully: %v\n", post)
	return nil
}

func (pr *PostRepository) Delete(id uint) error {
	log.Printf("Deleting post: %d\n", id)
	if err := pr.db.Delete(&model.Post{ID: id}).Error; err != nil {
		log.Printf("Error deleting post: %v\n", err)
		return err
	}
	if err := pr.s.Delete(id); err != nil {
		log.Printf("Error deleting post from search index: %v\n", err)
		return err
	}
	log.Printf("Post deleted successfully: %d\n", id)
	return nil
}

func (pr *PostRepository) Find(id uint) (*model.Post, error) {
	log.Printf("Finding post: %d\n", id)
	post := &model.Post{}
	if err := pr.db.Preload("Tags").Preload("Comments").Preload("Likes").First(post, id).Error; err != nil {
		log.Printf("Error finding post: %v\n", err)
		return nil, err
	}
	log.Printf("Post found: %v\n", post)
	return post, nil
}

func (pr *PostRepository) FindSome(skip, limit int) ([]*model.Post, error) {
	log.Printf("Finding some posts with skip: %d and limit: %d\n", skip, limit)
	var posts []*model.Post
	if limit < 0 {
		skip = -1
	}
	if err := pr.db.Preload("Tags").Preload("Comments").Preload("Likes").Offset(skip).Limit(limit).Find(&posts).Error; err != nil {
		log.Printf("Error finding some posts: %v\n", err)
		return nil, err
	}
	log.Printf("Found some posts: %v\n", posts)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PostedAt.After(posts[j].PostedAt)
	})
	return posts, nil
}

func (pr *PostRepository) FindByUser(userID uint, skip, limit int) ([]*model.Post, error) {
	log.Printf("Finding posts by user: %d with skip: %d and limit: %d\n", userID, skip, limit)
	var posts []*model.Post
	if limit < 0 {
		skip = -1
	}
	if err := pr.db.Preload("Tags").Preload("Comments").Preload("Likes").Where("user_id = ?", userID).Offset(skip).Limit(limit).Find(&posts).Error; err != nil {
		log.Printf("Error finding posts by user: %v\n", err)
		return nil, err
	}
	log.Printf("Found posts by user: %v\n", posts)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PostedAt.After(posts[j].PostedAt)
	})
	return posts, nil
}

func (pr *PostRepository) FindByTag(tagID uint, skip, limit int) ([]*model.Post, error) {
	log.Printf("Finding posts by tag: %d with skip: %d and limit: %d\n", tagID, skip, limit)
	var posts []*model.Post
	if limit < 0 {
		skip = -1
	}
	if err := pr.db.Preload("Tags").Preload("Comments").Preload("Likes").Joins("JOIN post_tags ON post_tags.post_id = posts.id").Where("post_tags.tag_id = ?", tagID).Offset(skip).Limit(limit).Find(&posts).Error; err != nil {
		log.Printf("Error finding posts by tag: %v\n", err)
		return nil, err
	}
	log.Printf("Found posts by tag: %v\n", posts)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PostedAt.After(posts[j].PostedAt)
	})
	return posts, nil
}

func (pr *PostRepository) ClearAssociations(postID uint) error {
	log.Printf("Clearing associations of post: %d\n", postID)
	if err := pr.db.Model(&model.Post{ID: postID}).Association("Tags").Clear(); err != nil {
		log.Printf("Error clearing associations: %v\n", err)
		return err
	}
	log.Printf("Associations cleared for post: %d\n", postID)
	return nil
}

func (pr *PostRepository) AddTag(postID, tagID uint) error {
	log.Printf("Adding tag: %d to post: %d\n", tagID, postID)
	if err := pr.db.Model(&model.Post{ID: postID}).Association("Tags").Append(&model.Tag{ID: tagID}); err != nil {
		log.Printf("Error adding tag: %v\n", err)
		return err
	}
	log.Printf("Tag added to post: %d\n", postID)
	return nil
}

func (pr *PostRepository) RemoveTag(postID, tagID uint) error {
	log.Printf("Removing tag: %d from post: %d\n", tagID, postID)
	if err := pr.db.Model(&model.Post{ID: postID}).Association("Tags").Delete(&model.Tag{ID: tagID}); err != nil {
		log.Printf("Error removing tag: %v\n", err)
		return err
	}
	log.Printf("Tag removed from post: %d\n", postID)
	return nil
}

func (pr *PostRepository) Search(query string, skip, limit int) ([]*model.Post, error) {
	log.Printf("Searching posts by query: %s with skip: %d and limit: %d\n", query, skip, limit)
	postIDs, err := pr.s.Search(query, skip, limit)
	if err != nil {
		log.Printf("Error searching posts: %v\n", err)
		return nil, err
	}

	var posts []*model.Post
	if err := pr.db.Preload("Tags").Preload("Comments").Preload("Likes").Where("id IN ?", postIDs).Find(&posts).Error; err != nil {
		log.Printf("Error finding posts by IDs: %v\n", err)
		return nil, err
	}
	log.Printf("Found posts by query: %v\n", posts)
	return posts, nil
}

func (pr *PostRepository) Statistic() (*model.PostStatistic, error) {
	log.Println("Getting post statistics")

	var stat model.PostStatistic
	err := pr.db.Model(&model.Post{}).
		Select("COUNT(*) AS total, " +
			"SUM(CASE WHEN edited_by != 0 THEN 1 ELSE 0 END) AS edited, " +
			"(SELECT COUNT(*) FROM likes WHERE deleted_at IS NULL) AS likes").
		Scan(&stat).Error

	if err != nil {
		log.Printf("Error getting post statistics: %v\n", err)
		return nil, err
	}

	log.Printf("Post statistics found: %+v\n", stat)
	return &stat, nil
}
