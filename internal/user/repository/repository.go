package repository

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
)

type Repository struct {
	User *UserRepository
	Ban  *BanRepository
}

func New(dialector gorm.Dialector) (*Repository, error) {
	log.Printf("Connecting %s...\n", dialector.Name())
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	log.Printf("Connected %s\n", dialector.Name())

	log.Println("Starting AutoMigrating...")
	if err := db.AutoMigrate(&model.User{}, &model.Ban{}); err != nil {
		return nil, err
	}
	log.Println("AutoMigrating completed")

	return &Repository{
		User: NewUserRepository(db),
		Ban:  NewBanRepository(db),
	}, nil
}
