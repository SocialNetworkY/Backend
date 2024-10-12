package service

import (
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"time"
)

type (
	BanRepo interface {
		Add(ban *model.Ban) error
		Save(ban *model.Ban) error
		Delete(id uint) error
		Find(id uint) (*model.Ban, error)
		FindByUser(userID uint, skip, limit int) ([]*model.Ban, error)
		FindByAdmin(adminID uint, skip, limit int) ([]*model.Ban, error)
	}

	BanService struct {
		repo BanRepo
	}
)

func NewBanService(r BanRepo) *BanService {
	return &BanService{
		repo: r,
	}
}

// Ban bans a user with time and reason
func (bs *BanService) Ban(userID, adminID uint, reason string, duration time.Duration) error {
	ban := &model.Ban{
		UserID:    userID,
		BannedBy:  adminID,
		BanReason: reason,
		BannedAt:  time.Now(),
		Duration:  duration,
		ExpiredAt: time.Now().Add(duration),
	}
	return bs.repo.Add(ban)
}

// Unban unbans a user with reason
func (bs *BanService) Unban(banID, adminID uint, reason string) error {
	ban, err := bs.repo.Find(banID)
	if err != nil {
		return err
	}

	ban.UnbanReason = reason
	ban.UnbannedBy = adminID
	ban.UnbannedAt = time.Now()
	return bs.repo.Save(ban)
}

// Find returns a ban by id
func (bs *BanService) Find(id uint) (*model.Ban, error) {
	return bs.repo.Find(id)
}

// FindByUser returns all bans for a user
func (bs *BanService) FindByUser(userID uint, skip, limit int) ([]*model.Ban, error) {
	return bs.repo.FindByUser(userID, skip, limit)
}
