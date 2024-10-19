package service

import (
	"context"
	"errors"
	"log"

	"github.com/SocialNetworkY/Backend/internal/report/model"
)

type (
	ReportRepo interface {
		Add(report *model.Report) error
		Save(report *model.Report) error
		Delete(id uint) error
		Get(id uint) (*model.Report, error)
		GetByPostUser(postID, userID uint) (*model.Report, error)
		GetSome(skip, limit int, status string) ([]*model.Report, error)
		GetByPost(postID uint, skip, limit int, status string) ([]*model.Report, error)
		GetByUser(userID uint, skip, limit int, status string) ([]*model.Report, error)
		GetByAdmin(adminID uint, skip, limit int, status string) ([]*model.Report, error)
		Search(query string, skip, limit int) ([]*model.Report, error)
	}

	PostGateway interface {
		PostInfo(ctx context.Context, postID uint) (*model.Post, error)
	}

	Report struct {
		repo ReportRepo
		pg   PostGateway
	}
)

func NewReport(repo ReportRepo, pg PostGateway) *Report {
	return &Report{
		repo: repo,
		pg:   pg,
	}
}

func (r *Report) Create(userID, postID uint, reason string) (*model.Report, error) {
	log.Printf("Creating report for user ID: %d and post ID: %d with reason: %s\n", userID, postID, reason)
	_, err := r.pg.PostInfo(context.Background(), postID)
	if err != nil {
		log.Printf("Error fetching post info: %v\n", err)
		return nil, errors.New("post not found")
	}

	_, err = r.repo.GetByPostUser(postID, userID)
	if err == nil {
		log.Printf("Report already exists for post ID: %d and user ID: %d\n", postID, userID)
		return nil, errors.New("report already exists")
	}

	report := &model.Report{
		UserID: userID,
		PostID: postID,
		Reason: reason,
		Status: model.ReportStatusPending,
	}

	if err := r.repo.Add(report); err != nil {
		log.Printf("Error adding report: %v\n", err)
		return nil, err
	}

	log.Printf("Report created successfully: %v\n", report)
	return report, nil
}

func (r *Report) Answer(reportID, adminID uint, answer string) (*model.Report, error) {
	log.Printf("Answering report ID: %d by admin ID: %d with answer: %s\n", reportID, adminID, answer)
	report, err := r.repo.Get(reportID)
	if err != nil {
		log.Printf("Error getting report: %v\n", err)
		return nil, err
	}

	report.AdminID = adminID
	report.Answer = answer
	report.Status = model.ReportStatusAnswered
	report.Closed = true

	if err := r.repo.Save(report); err != nil {
		log.Printf("Error saving report: %v\n", err)
		return nil, err
	}

	log.Printf("Report answered successfully: %v\n", report)
	return report, nil
}

func (r *Report) Reject(reportID, adminID uint, answer string) (*model.Report, error) {
	log.Printf("Rejecting report ID: %d by admin ID: %d with answer: %s\n", reportID, adminID, answer)
	report, err := r.repo.Get(reportID)
	if err != nil {
		log.Printf("Error getting report: %v\n", err)
		return nil, err
	}

	report.AdminID = adminID
	report.Answer = answer
	report.Status = model.ReportStatusRejected
	report.Closed = true

	if err := r.repo.Save(report); err != nil {
		log.Printf("Error saving report: %v\n", err)
		return nil, err
	}

	log.Printf("Report rejected successfully: %v\n", report)
	return report, nil
}

func (r *Report) Delete(reportID uint) error {
	log.Printf("Deleting report ID: %d\n", reportID)
	if err := r.repo.Delete(reportID); err != nil {
		log.Printf("Error deleting report: %v\n", err)
		return err
	}
	log.Printf("Report deleted successfully: %d\n", reportID)
	return nil
}

func (r *Report) DeleteByUser(userID uint) error {
	log.Printf("Deleting reports for user ID: %d\n", userID)
	reports, err := r.repo.GetByUser(userID, 0, -1, "")
	if err != nil {
		log.Printf("Error getting reports for user: %v\n", err)
		return err
	}

	for _, report := range reports {
		if err := r.Delete(report.ID); err != nil {
			log.Printf("Error deleting report ID: %d\n", report.ID)
			return err
		}
	}

	log.Printf("Reports deleted successfully for user ID: %d\n", userID)
	return nil
}

func (r *Report) DeleteByPost(postID uint) error {
	log.Printf("Deleting reports for post ID: %d\n", postID)
	reports, err := r.repo.GetByPost(postID, 0, -1, "")
	if err != nil {
		log.Printf("Error getting reports for post: %v\n", err)
		return err
	}

	for _, report := range reports {
		if err := r.Delete(report.ID); err != nil {
			log.Printf("Error deleting report ID: %d\n", report.ID)
			return err
		}
	}

	log.Printf("Reports deleted successfully for post ID: %d\n", postID)
	return nil
}

func (r *Report) Get(reportID uint) (*model.Report, error) {
	log.Printf("Getting report ID: %d\n", reportID)
	report, err := r.repo.Get(reportID)
	if err != nil {
		log.Printf("Error getting report: %v\n", err)
		return nil, err
	}
	log.Printf("Report found: %v\n", report)
	return report, nil
}

func (r *Report) GetSome(skip, limit int, status string) ([]*model.Report, error) {
	log.Printf("Getting some reports with skip: %d, limit: %d, and status: %s\n", skip, limit, status)
	reports, err := r.repo.GetSome(skip, limit, status)
	if err != nil {
		log.Printf("Error getting some reports: %v\n", err)
		return nil, err
	}
	log.Printf("Reports found: %v\n", reports)
	return reports, nil
}

func (r *Report) GetByPost(postID uint, skip, limit int, status string) ([]*model.Report, error) {
	log.Printf("Getting reports for post ID: %d with skip: %d, limit: %d, and status: %s\n", postID, skip, limit, status)
	reports, err := r.repo.GetByPost(postID, skip, limit, status)
	if err != nil {
		log.Printf("Error getting reports for post: %v\n", err)
		return nil, err
	}
	log.Printf("Reports found: %v\n", reports)
	return reports, nil
}

func (r *Report) GetByUser(userID uint, skip, limit int, status string) ([]*model.Report, error) {
	log.Printf("Getting reports for user ID: %d with skip: %d, limit: %d, and status: %s\n", userID, skip, limit, status)
	reports, err := r.repo.GetByUser(userID, skip, limit, status)
	if err != nil {
		log.Printf("Error getting reports for user: %v\n", err)
		return nil, err
	}
	log.Printf("Reports found: %v\n", reports)
	return reports, nil
}

func (r *Report) GetByAdmin(adminID uint, skip, limit int, status string) ([]*model.Report, error) {
	log.Printf("Getting reports for admin ID: %d with skip: %d, limit: %d, and status: %s\n", adminID, skip, limit, status)
	reports, err := r.repo.GetByAdmin(adminID, skip, limit, status)
	if err != nil {
		log.Printf("Error getting reports for admin: %v\n", err)
		return nil, err
	}
	log.Printf("Reports found: %v\n", reports)
	return reports, nil
}

func (r *Report) Search(query string, skip, limit int) ([]*model.Report, error) {
	log.Printf("Searching reports with query: %s, skip: %d and limit: %d\n", query, skip, limit)
	reports, err := r.repo.Search(query, skip, limit)
	if err != nil {
		log.Printf("Error searching reports: %v\n", err)
		return nil, err
	}
	log.Printf("Reports found: %v\n", reports)
	return reports, nil
}
