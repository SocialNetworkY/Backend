package repository

import (
	"github.com/SocialNetworkY/Backend/internal/report/model"
	"gorm.io/gorm"
)

type Report struct {
	db *gorm.DB
}

func NewReport(db *gorm.DB) *Report {
	return &Report{
		db: db,
	}
}

func (r *Report) Add(report *model.Report) error {
	return r.db.Create(report).Error
}

func (r *Report) Save(report *model.Report) error {
	return r.db.Save(report).Error
}

func (r *Report) Delete(id uint) error {
	return r.db.Delete(&model.Report{}, id).Error
}

func (r *Report) Get(id uint) (*model.Report, error) {
	var report model.Report
	err := r.db.First(&report, id).Error
	return &report, err
}

func (r *Report) GetByPostUser(postID, userID uint) (*model.Report, error) {
	var report model.Report
	err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&report).Error
	return &report, err
}

func (r *Report) GetSome(skip, limit int, status string) ([]*model.Report, error) {
	var reports []*model.Report
	if limit < 0 {
		skip = 0
	}
	query := r.db.Offset(skip).Limit(limit)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&reports).Error
	return reports, err
}

func (r *Report) GetByPost(postID uint, skip, limit int, status string) ([]*model.Report, error) {
	var reports []*model.Report
	if limit < 0 {
		skip = 0
	}
	query := r.db.Where("post_id = ?", postID).Offset(skip).Limit(limit)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&reports).Error
	return reports, err
}

func (r *Report) GetByUser(userID uint, skip, limit int, status string) ([]*model.Report, error) {
	var reports []*model.Report
	if limit < 0 {
		skip = 0
	}
	query := r.db.Where("user_id = ?", userID).Offset(skip).Limit(limit)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&reports).Error
	return reports, err
}

func (r *Report) GetByAdmin(adminID uint, skip, limit int, status string) ([]*model.Report, error) {
	var reports []*model.Report
	if limit < 0 {
		skip = 0
	}
	query := r.db.Where("admin_id", adminID).Offset(skip).Limit(limit)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&reports).Error
	return reports, err
}
