package model

import (
	"gorm.io/gorm"
	"time"
)

type Ban struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	UserID      uint          `json:"userID"`
	BannedBy    uint          `json:"bannedBy"`
	BanReason   string        `json:"banReason"`
	BannedAt    time.Time     `json:"bannedAt"`
	Duration    time.Duration `json:"duration"`
	ExpiredAt   time.Time     `json:"expiredAt"`
	UnbanReason string        `json:"unbanReason"`
	UnbannedBy  uint          `json:"unbannedBy"`
	UnbannedAt  time.Time     `json:"unbannedAt" gorm:"default:null"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Active bool ` json:"active" gorm:"-"`
}

func (b *Ban) AfterFind(tx *gorm.DB) (err error) {
	b.Active = b.ExpiredAt.After(time.Now()) && b.UnbanReason == ""
	return nil
}
