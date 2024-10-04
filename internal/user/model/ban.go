package model

import (
	"gorm.io/gorm"
	"time"
)

type Ban struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	UserID      uint          `json:"userID"`
	BannedBy    uint          `json:"bannedBy"`
	BanReason   string        `json:"banReason"`
	BannedAt    time.Time     `json:"bannedAt"`
	Duration    time.Duration `json:"duration"`
	ExpiredAt   time.Time     `json:"expiredAt"`
	UnbanReason string        `json:"unbanReason"`
	UnbannedBy  uint          `json:"unbannedBy"`
	UnbannedAt  time.Time     `json:"unbannedAt"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Active bool `gorm:"-" json:"active"`
}

func (b *Ban) BeforeLoad() {
	b.Active = b.ExpiredAt.After(time.Now()) && b.UnbanReason == ""
}
