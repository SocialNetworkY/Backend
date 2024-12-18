package service

import (
	"context"
	"time"

	"github.com/SocialNetworkY/Backend/internal/user/model"
)

type (
	BanRepo interface {
		Add(ban *model.Ban) error
		Save(ban *model.Ban) error
		Delete(id uint) error
		Find(id uint) (*model.Ban, error)
		FindSome(skip, limit int) ([]*model.Ban, error)
		FindByUser(userID uint, skip, limit int) ([]*model.Ban, error)
		FindByAdmin(adminID uint, skip, limit int) ([]*model.Ban, error)
		Search(query string, skip, limit int) ([]*model.Ban, error)
		Statistic() (*model.BanStatistic, error)
	}

	BanService struct {
		repo BanRepo
		pg   PostGateway
	}
)

func NewBanService(r BanRepo, pg PostGateway) *BanService {
	return &BanService{
		repo: r,
		pg:   pg,
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

	if err := bs.pg.DeleteUserPosts(context.Background(), userID); err != nil {
		return err
	}

	if err := bs.pg.DeleteUserComments(context.Background(), userID); err != nil {
		return err
	}

	if err := bs.pg.DeleteUserLikes(context.Background(), userID); err != nil {
		return err
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

// FindSome returns some bans
func (bs *BanService) FindSome(skip, limit int) ([]*model.Ban, error) {
	return bs.repo.FindSome(skip, limit)
}

// Search returns bans by query
func (bs *BanService) Search(query string, skip, limit int) ([]*model.Ban, error) {
	return bs.repo.Search(query, skip, limit)
}

func (bs *BanService) Statistic() (*model.BanStatistic, error) {
	return bs.repo.Statistic()
}
