package mysql

import (
	"github.com/lapkomo2018/goTwitterAuthService/pkg/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type Storage struct {
	User            *UserStorage
	RefreshToken    *RefreshTokenStorage
	ActivationToken *ActivationTokenStorage
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
	if err := db.AutoMigrate(&model.User{}, &model.ActivationToken{}, &model.RefreshToken{}); err != nil {
		return nil, err
	}
	log.Println("AutoMigrating completed")

	return &Storage{
		User:            NewUserStorage(db),
		RefreshToken:    NewRefreshTokenStorage(db),
		ActivationToken: NewActivationTokenStorage(db),
	}, nil
}
