package model

import (
	"gorm.io/gorm"
	"time"
)

type Ban struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	AdminID     uint
	Reason      string
	Duration    time.Duration
	ExpiredAt   time.Time
	UnbanReason string
	UnbannedBy  uint
	UnbannedAt  time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Active bool `gorm:"-"`
}

func (b *Ban) BeforeLoad() {
	b.Active = b.ExpiredAt.After(time.Now()) && b.UnbanReason == ""
}
