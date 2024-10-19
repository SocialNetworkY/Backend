package repository

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/SocialNetworkY/Backend/internal/report/model"
)

type Repository struct {
	Report *Report
}

func New(dialector gorm.Dialector, rs ReportSearch) (*Repository, error) {
	log.Printf("Connecting %s...\n", dialector.Name())
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	log.Printf("Connected %s\n", dialector.Name())

	log.Println("Starting AutoMigrating...")
	if err := db.AutoMigrate(&model.Report{}); err != nil {
		return nil, err
	}
	log.Println("AutoMigrating completed")

	return &Repository{
		Report: NewReport(db, rs),
	}, nil
}
