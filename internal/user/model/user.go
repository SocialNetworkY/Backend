package model

import (
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

const (
	RoleUser      = 0
	RoleAdminLvl1 = 1
	RoleAdminLvl2 = 2
	RoleAdminLvl3 = 3
)

func (u *User) CheckBans() {
	u.Banned = false
	for _, ban := range u.Bans {
		if ban.ExpiredAt.After(time.Now()) {
			u.Banned = true
		}
	}
}

func (u *User) BeforeLoad(tx *gorm.DB) (err error) {
	if err := tx.Preload("Bans").Find(u).Error; err != nil {
		return err
	}
	u.CheckBans()
	u.Admin = u.Role > 0
	return
}
