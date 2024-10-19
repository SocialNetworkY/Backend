package repository

import (
	"log"

	"github.com/SocialNetworkY/Backend/internal/post/model"
	"gorm.io/gorm"
)

type (
	TagSearch interface {
		Index(tag *model.Tag) error
		Delete(id uint) error
		Search(query string, skip, limit int) ([]uint, error)
	}
	TagRepository struct {
		db *gorm.DB
		s  TagSearch
	}
)

func NewTagRepository(db *gorm.DB, s TagSearch) *TagRepository {
	return &TagRepository{
		db: db,
		s:  s,
	}
}

func (tr *TagRepository) Add(tag *model.Tag) error {
	log.Printf("Adding tag: %v\n", tag)
	if err := tr.db.Create(tag).Error; err != nil {
		log.Printf("Error adding tag: %v\n", err)
		return err
	}
	if err := tr.s.Index(tag); err != nil {
		log.Printf("Error indexing tag: %v\n", err)
		return err
	}
	log.Printf("Tag added successfully: %v\n", tag)
	return nil
}

func (tr *TagRepository) Save(tag *model.Tag) error {
	log.Printf("Saving tag: %v\n", tag)
	if err := tr.db.Save(tag).Error; err != nil {
		log.Printf("Error saving tag: %v\n", err)
		return err
	}
	if err := tr.s.Index(tag); err != nil {
		log.Printf("Error indexing tag: %v\n", err)
		return err
	}
	log.Printf("Tag saved successfully: %v\n", tag)
	return nil
}

func (tr *TagRepository) Delete(id uint) error {
	log.Printf("Deleting tag by id: %d\n", id)
	if err := tr.db.Delete(&model.Tag{ID: id}).Error; err != nil {
		log.Printf("Error deleting tag: %v\n", err)
		return err
	}
	if err := tr.s.Delete(id); err != nil {
		log.Printf("Error deleting tag from search index: %v\n", err)
		return err
	}
	log.Printf("Tag deleted successfully: %d\n", id)
	return nil
}

func (tr *TagRepository) Find(id uint) (*model.Tag, error) {
	log.Printf("Finding tag by id: %d\n", id)
	tag := &model.Tag{}
	if err := tr.db.Preload("Posts").First(tag, id).Error; err != nil {
		log.Printf("Error finding tag: %v\n", err)
		return nil, err
	}
	log.Printf("Tag found: %v\n", tag)
	return tag, nil
}

func (tr *TagRepository) FindByName(name string) (*model.Tag, error) {
	log.Printf("Finding tag by name: %s\n", name)
	tag := &model.Tag{}
	if err := tr.db.Preload("Posts").Where("name = ?", name).First(tag).Error; err != nil {
		log.Printf("Error finding tag: %v\n", err)
		return nil, err
	}
	log.Printf("Tag found: %v\n", tag)
	return tag, nil
}

func (tr *TagRepository) FindSome(skip, limit int) ([]*model.Tag, error) {
	log.Printf("Finding some tags with skip: %d and limit: %d\n", skip, limit)
	var tags []*model.Tag
	if limit < 0 {
		skip = -1
	}
	if err := tr.db.Preload("Posts").Offset(skip).Limit(limit).Find(&tags).Error; err != nil {
		log.Printf("Error finding some tags: %v\n", err)
		return nil, err
	}
	log.Printf("Found some tags: %v\n", tags)
	return tags, nil
}

func (tr *TagRepository) FindByPost(postID uint, skip, limit int) ([]*model.Tag, error) {
	log.Printf("Finding tags by post: %d\n", postID)
	var tags []*model.Tag
	if limit < 0 {
		skip = -1
	}
	if err := tr.db.Preload("Posts").Joins("JOIN post_tags ON post_tags.tag_id = tags.id").Where("post_tags.post_id = ?", postID).Offset(skip).Limit(limit).Find(&tags).Error; err != nil {
		log.Printf("Error finding tags: %v\n", err)
		return nil, err
	}
	log.Printf("Tags found: %v\n", tags)
	return tags, nil
}

func (tr *TagRepository) ClearAssociations(tagID uint) error {
	log.Printf("Clearing associations of tag: %d\n", tagID)
	if err := tr.db.Model(&model.Tag{ID: tagID}).Association("Posts").Clear(); err != nil {
		log.Printf("Error clearing associations: %v\n", err)
		return err
	}
	log.Printf("Associations cleared for tag: %d\n", tagID)
	return nil
}

func (tr *TagRepository) AddPost(tagID, postID uint) error {
	log.Printf("Adding post: %d to tag: %d\n", postID, tagID)
	if err := tr.db.Model(&model.Tag{ID: tagID}).Association("Posts").Append(&model.Post{ID: postID}); err != nil {
		log.Printf("Error adding post to tag: %v\n", err)
		return err
	}
	log.Printf("Post: %d added to tag: %d\n", postID, tagID)
	return nil
}

func (tr *TagRepository) RemovePost(tagID, postID uint) error {
	log.Printf("Removing post: %d from tag: %d\n", postID, tagID)
	if err := tr.db.Model(&model.Tag{ID: tagID}).Association("Posts").Delete(&model.Post{ID: postID}); err != nil {
		log.Printf("Error removing post from tag: %v\n", err)
		return err
	}
	log.Printf("Post: %d removed from tag: %d\n", postID, tagID)
	return nil
}

func (tr *TagRepository) Search(query string, skip, limit int) ([]*model.Tag, error) {
	log.Printf("Searching tags by query: %s with skip: %d and limit: %d\n", query, skip, limit)
	ids, err := tr.s.Search(query, skip, limit)
	if err != nil {
		log.Printf("Error searching tags: %v\n", err)
		return nil, err
	}

	var tags []*model.Tag
	if err := tr.db.Preload("Posts").Where("id IN ?", ids).Find(&tags).Error; err != nil {
		log.Printf("Error finding tags by IDs: %v\n", err)
		return nil, err
	}
	log.Printf("Tags found: %v\n", tags)
	return tags, nil
}
