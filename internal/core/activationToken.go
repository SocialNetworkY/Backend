package core

type ActivationToken struct {
	UserID uint   `gorm:"primaryKey"`
	Token  string `gorm:"unique"`
}
