package model

type RefreshToken struct {
	UserID uint   `gorm:"primaryKey"`
	Token  string `gorm:"unique"`
}
