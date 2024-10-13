package repository

import (
	"github.com/SocialNetworkY/Backend/internal/auth/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type Repository struct {
	User            *UserRepository
	RefreshToken    *RefreshTokenRepository
	ActivationToken *ActivationTokenRepository
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
	if err := db.AutoMigrate(&model.User{}, &model.ActivationToken{}, &model.RefreshToken{}); err != nil {
		return nil, err
	}
	log.Println("AutoMigrating completed")

	return &Repository{
		User:            NewUserRepository(db),
		RefreshToken:    NewRefreshTokenRepository(db),
		ActivationToken: NewActivationTokenRepository(db),
	}, nil
}
