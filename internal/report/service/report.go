package service

import (
	"context"
	"errors"

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
	_, err := r.pg.PostInfo(context.Background(), postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	_, err = r.repo.GetByPostUser(postID, userID)
	if err == nil {
		return nil, errors.New("report already exists")
	}

	report := &model.Report{
		UserID: userID,
		PostID: postID,
		Reason: reason,
		Status: model.ReportStatusPending,
	}

	if err := r.repo.Add(report); err != nil {
		return nil, err
	}

	return report, nil
}

func (r *Report) Answer(reportID, adminID uint, answer string) (*model.Report, error) {
	report, err := r.repo.Get(reportID)
	if err != nil {
		return nil, err
	}

	report.AdminID = adminID
	report.Answer = answer
	report.Status = model.ReportStatusAnswered
	report.Closed = true

	if err := r.repo.Save(report); err != nil {
		return nil, err
	}

	return report, nil
}

func (r *Report) Reject(reportID, adminID uint, answer string) (*model.Report, error) {
	report, err := r.repo.Get(reportID)
	if err != nil {
		return nil, err
	}

	report.AdminID = adminID
	report.Answer = answer
	report.Status = model.ReportStatusRejected
	report.Closed = true

	if err := r.repo.Save(report); err != nil {
		return nil, err
	}

	return report, nil
}

func (r *Report) Delete(reportID uint) error {
	return r.repo.Delete(reportID)
}

func (r *Report) DeleteByUser(userID uint) error {
	reports, err := r.repo.GetByUser(userID, 0, -1, "")
	if err != nil {
		return err
	}

	for _, report := range reports {
		if err := r.Delete(report.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r *Report) DeleteByPost(postID uint) error {
	reports, err := r.repo.GetByPost(postID, 0, -1, "")
	if err != nil {
		return err
	}

	for _, report := range reports {
		if err := r.Delete(report.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r *Report) Get(reportID uint) (*model.Report, error) {
	return r.repo.Get(reportID)
}

func (r *Report) GetSome(skip, limit int, status string) ([]*model.Report, error) {
	return r.repo.GetSome(skip, limit, status)
}

func (r *Report) GetByPost(postID uint, skip, limit int, status string) ([]*model.Report, error) {
	return r.repo.GetByPost(postID, skip, limit, status)
}

func (r *Report) GetByUser(userID uint, skip, limit int, status string) ([]*model.Report, error) {
	return r.repo.GetByUser(userID, skip, limit, status)
}

func (r *Report) GetByAdmin(adminID uint, skip, limit int, status string) ([]*model.Report, error) {
	return r.repo.GetByAdmin(adminID, skip, limit, status)
}
