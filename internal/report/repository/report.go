package repository

import (
	"log"

	"github.com/SocialNetworkY/Backend/internal/report/model"
	"gorm.io/gorm"
)

type (
	ReportSearch interface {
		Index(report *model.Report) error
		Delete(id uint) error
		Search(query string, skip, limit int) ([]uint, error)
	}

	Report struct {
		db *gorm.DB
		s  ReportSearch
	}
)

func NewReport(db *gorm.DB, s ReportSearch) *Report {
	return &Report{
		db: db,
		s:  s,
	}
}

func (r *Report) Add(report *model.Report) error {
	log.Printf("Adding report: %v\n", report)
	if err := r.db.Create(report).Error; err != nil {
		log.Printf("Error adding report: %v\n", err)
		return err
	}

	if err := r.s.Index(report); err != nil {
		log.Printf("Error indexing report: %v\n", err)
		return err
	}
	log.Printf("Report added successfully: %v\n", report)
	return nil
}

func (r *Report) Save(report *model.Report) error {
	log.Printf("Saving report: %v\n", report)
	if err := r.db.Save(report).Error; err != nil {
		log.Printf("Error saving report: %v\n", err)
		return err
	}
	if err := r.s.Index(report); err != nil {
		log.Printf("Error indexing report: %v\n", err)
		return err
	}
	log.Printf("Report saved successfully: %v\n", report)
	return nil
}

func (r *Report) Delete(id uint) error {
	log.Printf("Deleting report: %d\n", id)
	if err := r.db.Delete(&model.Report{ID: id}).Error; err != nil {
		log.Printf("Error deleting report: %v\n", err)
		return err
	}
	if err := r.s.Delete(id); err != nil {
		log.Printf("Error deleting report from search index: %v\n", err)
		return err
	}
	log.Printf("Report deleted successfully: %d\n", id)
	return nil
}

func (r *Report) Get(id uint) (*model.Report, error) {
	log.Printf("Getting report with ID: %d\n", id)
	var report model.Report
	err := r.db.First(&report, id).Error
	if err != nil {
		log.Printf("Error getting report: %v\n", err)
		return nil, err
	}
	log.Printf("Report found: %v\n", report)
	return &report, nil
}

func (r *Report) GetByPostUser(postID, userID uint) (*model.Report, error) {
	log.Printf("Getting report for post ID: %d and user ID: %d\n", postID, userID)
	var report model.Report
	err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&report).Error
	if err != nil {
		log.Printf("Error getting report: %v\n", err)
		return nil, err
	}
	log.Printf("Report found: %v\n", report)
	return &report, nil
}

func (r *Report) GetSome(skip, limit int, status string) ([]*model.Report, error) {
	log.Printf("Getting some reports with skip: %d, limit: %d, and status: %s\n", skip, limit, status)
	var reports []*model.Report
	if limit < 0 {
		skip = 0
	}
	query := r.db.Offset(skip).Limit(limit)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&reports).Error
	if err != nil {
		log.Printf("Error getting some reports: %v\n", err)
		return nil, err
	}
	log.Printf("Reports found: %v\n", reports)
	return reports, nil
}

func (r *Report) GetByPost(postID uint, skip, limit int, status string) ([]*model.Report, error) {
	log.Printf("Getting reports for post ID: %d with skip: %d, limit: %d, and status: %s\n", postID, skip, limit, status)
	var reports []*model.Report
	if limit < 0 {
		skip = 0
	}
	query := r.db.Where("post_id = ?", postID).Offset(skip).Limit(limit)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&reports).Error
	if err != nil {
		log.Printf("Error getting reports for post: %v\n", err)
		return nil, err
	}
	log.Printf("Reports found: %v\n", reports)
	return reports, nil
}

func (r *Report) GetByUser(userID uint, skip, limit int, status string) ([]*model.Report, error) {
	log.Printf("Getting reports for user ID: %d with skip: %d, limit: %d, and status: %s\n", userID, skip, limit, status)
	var reports []*model.Report
	if limit < 0 {
		skip = 0
	}
	query := r.db.Where("user_id = ?", userID).Offset(skip).Limit(limit)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&reports).Error
	if err != nil {
		log.Printf("Error getting reports for user: %v\n", err)
		return nil, err
	}
	log.Printf("Reports found: %v\n", reports)
	return reports, nil
}

func (r *Report) GetByAdmin(adminID uint, skip, limit int, status string) ([]*model.Report, error) {
	log.Printf("Getting reports for admin ID: %d with skip: %d, limit: %d, and status: %s\n", adminID, skip, limit, status)
	var reports []*model.Report
	if limit < 0 {
		skip = 0
	}
	query := r.db.Where("admin_id", adminID).Offset(skip).Limit(limit)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&reports).Error
	if err != nil {
		log.Printf("Error getting reports for admin: %v\n", err)
		return nil, err
	}
	log.Printf("Reports found: %v\n", reports)
	return reports, nil
}

func (r *Report) Search(query string, skip, limit int) ([]*model.Report, error) {
	log.Printf("Searching reports with query: %s, skip: %d and limit: %d\n", query, skip, limit)
	ids, err := r.s.Search(query, skip, limit)
	if err != nil {
		log.Printf("Error searching reports: %v\n", err)
		return nil, err
	}

	var reports []*model.Report
	if err := r.db.Where("id IN ?", ids).Find(&reports).Error; err != nil {
		log.Printf("Error getting reports: %v\n", err)
		return nil, err
	}
	log.Printf("Found reports: %v\n", reports)
	return reports, nil
}

func (r *Report) Statistic() (*model.ReportStatistic, error) {
	log.Println("Getting report statistic")

	var stat model.ReportStatistic
	err := r.db.Model(&Report{}).
		Select("COUNT(*) AS total, "+
			"SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS pending, "+
			"SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS answered, "+
			"SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS rejected",
			model.ReportStatusPending, model.ReportStatusAnswered, model.ReportStatusRejected).
		Scan(&stat).Error

	if err != nil {
		log.Printf("Error getting report statistic: %v\n", err)
		return nil, err
	}

	log.Printf("Report statistic found: %+v\n", stat)
	return &stat, nil
}
