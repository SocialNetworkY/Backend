package repository

import (
	"github.com/SocialNetworkY/Backend/internal/user/model"
	"gorm.io/gorm"
)

type (
	BanSearch interface {
		Index(ban *model.Ban) error
		Delete(id uint) error
		Search(query string, skip, limit int) ([]uint, error)
	}

	BanRepository struct {
		db *gorm.DB
		s  BanSearch
	}
)

func NewBanRepository(db *gorm.DB, s BanSearch) *BanRepository {
	return &BanRepository{
		db: db,
		s:  s,
	}
}

func (br *BanRepository) Add(ban *model.Ban) error {
	if err := br.db.Create(ban).Error; err != nil {
		return err
	}
	return br.s.Index(ban)
}

func (br *BanRepository) Save(ban *model.Ban) error {
	if err := br.db.Save(ban).Error; err != nil {
		return err
	}
	return br.s.Index(ban)
}

func (br *BanRepository) Delete(id uint) error {
	if err := br.db.Delete(&model.Ban{}, id).Error; err != nil {
		return err
	}
	return br.s.Delete(id)
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
	if limit < 0 {
		skip = -1
	}
	if err := br.db.Offset(skip).Limit(limit).Find(&bans).Error; err != nil {
		return nil, err
	}
	return bans, nil
}

func (br *BanRepository) FindByUser(userID uint, skip, limit int) ([]*model.Ban, error) {
	var bans []*model.Ban
	if limit < 0 {
		skip = -1
	}
	if err := br.db.Offset(skip).Where("user_id = ?", userID).Limit(limit).Find(&bans).Error; err != nil {
		return nil, err
	}
	return bans, nil
}

func (br *BanRepository) FindByAdmin(adminID uint, skip, limit int) ([]*model.Ban, error) {
	var bans []*model.Ban
	if limit < 0 {
		skip = -1
	}
	if err := br.db.Offset(skip).Where("admin_id = ?", adminID).Limit(limit).Find(&bans).Error; err != nil {
		return nil, err
	}
	return bans, nil
}

func (br *BanRepository) Search(query string, skip, limit int) ([]*model.Ban, error) {
	bansIDs, err := br.s.Search(query, skip, limit)
	if err != nil {
		return nil, err
	}

	var bans []*model.Ban
	if err := br.db.Where("id IN ?", bansIDs).Find(&bans).Error; err != nil {
		return nil, err
	}

	return bans, nil
}
