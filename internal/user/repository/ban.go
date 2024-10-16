package repository

import (
	"github.com/SocialNetworkY/Backend/internal/user/model"
	"gorm.io/gorm"
)

type BanRepository struct {
	db *gorm.DB
}

func NewBanRepository(db *gorm.DB) *BanRepository {
	return &BanRepository{
		db: db,
	}
}

func (br *BanRepository) Add(ban *model.Ban) error {
	if err := br.db.Create(ban).Error; err != nil {
		return err
	}
	return nil
}

func (br *BanRepository) Save(ban *model.Ban) error {
	if err := br.db.Save(ban).Error; err != nil {
		return err
	}
	return nil
}

func (br *BanRepository) Delete(id uint) error {
	if err := br.db.Delete(&model.Ban{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (br *BanRepository) Find(id uint) (*model.Ban, error) {
	ban := &model.Ban{}
	if err := br.db.First(ban, id).Error; err != nil {
		return nil, err
	}
	return ban, nil
}

func (br *BanRepository) FindSome(skip, limit int) ([]*model.Ban, error) {
	var bans []*model.Ban
	query := br.db.Offset(skip)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&bans).Error; err != nil {
		return nil, err
	}
	return bans, nil
}

func (br *BanRepository) FindByUser(userID uint, skip, limit int) ([]*model.Ban, error) {
	var bans []*model.Ban
	query := br.db.Offset(skip).Where("user_id = ?", userID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&bans).Error; err != nil {
		return nil, err
	}
	return bans, nil
}

func (br *BanRepository) FindByAdmin(adminID uint, skip, limit int) ([]*model.Ban, error) {
	var bans []*model.Ban
	query := br.db.Offset(skip).Where("admin_id = ?", adminID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&bans).Error; err != nil {
		return nil, err
	}
	return bans, nil
}
