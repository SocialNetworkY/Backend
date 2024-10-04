package service

import (
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"time"
)

type (
	BanStorage interface {
		Add(ban *model.Ban) error
		Save(ban *model.Ban) error
		Delete(id uint) error
		Find(id uint) (*model.Ban, error)
		FindByUserID(userID uint) ([]*model.Ban, error)
		FindByAdminID(adminID uint) ([]*model.Ban, error)
	}

	BanService struct {
		s BanStorage
	}
)

func NewBanService(s BanStorage) *BanService {
	return &BanService{
		s: s,
	}
}

// BanUser bans a user with time and reason
func (bs *BanService) BanUser(userID, adminID uint, reason string, duration time.Duration) error {
	ban := &model.Ban{
		UserID:    userID,
		BannedBy:  adminID,
		BanReason: reason,
		BannedAt:  time.Now(),
		Duration:  duration,
		ExpiredAt: time.Now().Add(duration),
	}
	return bs.s.Add(ban)
}

// UnbanByBanID unbans a user with reason
func (bs *BanService) UnbanByBanID(banID, adminID uint, reason string) error {
	ban, err := bs.s.Find(banID)
	if err != nil {
		return err
	}

	ban.UnbanReason = reason
	ban.UnbannedBy = adminID
	ban.UnbannedAt = time.Now()
	return bs.s.Save(ban)
}

// FindBan returns a ban by id
func (bs *BanService) FindBan(id uint) (*model.Ban, error) {
	return bs.s.Find(id)
}

// FindBansByUserID returns all bans for a user
func (bs *BanService) FindBansByUserID(userID uint) ([]*model.Ban, error) {
	return bs.s.FindByUserID(userID)
}
