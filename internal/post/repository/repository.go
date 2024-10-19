package repository

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/SocialNetworkY/Backend/internal/post/model"
)

type Repository struct {
	Post    *PostRepository
	Comment *CommentRepository
	Like    *LikeRepository
	Tag     *TagRepository
}

func New(dialector gorm.Dialector, ps PostSearch, cs CommentSearch, ts TagSearch) (*Repository, error) {
	log.Printf("Connecting %s...\n", dialector.Name())
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	log.Printf("Connected %s\n", dialector.Name())

	log.Println("Starting AutoMigrating...")
	if err := db.AutoMigrate(&model.Post{}, &model.Like{}, &model.Comment{}, &model.Tag{}); err != nil {
		return nil, err
	}
	log.Println("AutoMigrating completed")

	return &Repository{
		Post:    NewPostRepository(db, ps),
		Comment: NewCommentRepository(db, cs),
		Like:    NewLikeRepository(db),
		Tag:     NewTagRepository(db, ts),
	}, nil
}
