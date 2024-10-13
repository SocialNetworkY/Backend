package model

import (
	"github.com/SocialNetworkY/Backend/pkg/constant"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"unique"`
	Username  string         `json:"username" gorm:"unique"`
	Nickname  string         `json:"nickname"`
	Avatar    string         `json:"avatar"`
	Role      uint           `json:"role"`
	Banned    bool           `json:"banned" gorm:"-"`
	Admin     bool           `json:"admin" gorm:"-"`
	ActiveBan *Ban           `json:"activeBan" gorm:"-"`
	Bans      []*Ban         `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// AfterFind is a gorm hook that is called after a find operation
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	for _, ban := range u.Bans {
		if ban.Active {
			u.Banned = true
			u.ActiveBan = ban
		}
	}

	u.Admin = u.Role > constant.RoleUser
	return nil
}
