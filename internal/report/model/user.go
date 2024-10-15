package model

import "time"

type User struct {
	ID           uint
	Role         uint
	Banned       bool
	BanReason    string
	BanExpiredAt time.Time
}
