package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	ReportStatusPending  = "pending"
	ReportStatusAnswered = "answered"
	ReportStatusRejected = "rejected"
)

type Report struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id"`
	PostID    uint           `json:"post_id"`
	Reason    string         `json:"reason"`
	AdminID   uint           `json:"admin_id"`
	Answer    string         `json:"answer"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Closed bool `json:"closed" gorm:"-"`
}

func (r *Report) AfterFind(tx *gorm.DB) (err error) {
	if r.Status != ReportStatusPending {
		r.Closed = true
	}
	return nil
}

type ReportStatistic struct {
	Total    int64 `json:"total"`
	Pending  int64 `json:"pending"`
	Answered int64 `json:"answered"`
	Rejected int64 `json:"rejected"`
}
