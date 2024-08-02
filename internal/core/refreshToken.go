package core

import "time"

type RefreshToken struct {
	Token     string `gorm:"primaryKey"`
	UserID    uint
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
}
