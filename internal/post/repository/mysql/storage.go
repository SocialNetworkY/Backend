package mysql

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
)

type Storage struct {
	Post    *PostStorage
	Comment *CommentStorage
	Like    *LikeStorage
	Tag     *TagStorage
}

func New(dsn string) (*Storage, error) {
	log.Println("Connecting mysql...")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	log.Println("Connected mysql")

	log.Println("Starting AutoMigrating...")
	if err := db.AutoMigrate(&model.Post{}, &model.Like{}, &model.Comment{}, &model.Tag{}); err != nil {
		return nil, err
	}
	log.Println("AutoMigrating completed")

	return &Storage{
		Post:    NewPostStorage(db),
		Comment: NewCommentStorage(db),
		Like:    NewLikeStorage(db),
		Tag:     NewTagStorage(db),
	}, nil
}
