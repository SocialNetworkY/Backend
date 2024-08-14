package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Email       string `gorm:"unique"`
	Username    string `gorm:"unique"`
	Password    string
	IsActivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	if err := tx.Where("user_id = ?", u.ID).Delete(&RefreshToken{}).Error; err != nil {
		return ErrUserRefreshTokenDelete
	}
	if err := tx.Where("user_id = ?", u.ID).Delete(&ActivationToken{}).Error; err != nil {
		return ErrUserActivationTokenDelete
	}

	return nil
}
