package model

import (
	"github.com/lapkomo2018/goTwitterServices/pkg/constant"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique"`
	Username  string `gorm:"unique"`
	Nickname  string
	Avatar    string
	Role      uint
	Banned    bool  `gorm:"-"`
	Admin     bool  `gorm:"-"`
	Bans      []Ban `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) CheckBans() {
	u.Banned = false
	for _, ban := range u.Bans {
		if ban.Active {
			u.Banned = true
		}
	}
}

func (u *User) BeforeLoad(tx *gorm.DB) (err error) {
	if err := tx.Preload("Bans").Find(u).Error; err != nil {
		return err
	}
	u.CheckBans()
	u.Admin = u.Role > constant.RoleUser
	return
}
