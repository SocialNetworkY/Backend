package mysql

import (
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"gorm.io/gorm"
)

type BanStorage struct {
	db *gorm.DB
}

func NewBanStorage(db *gorm.DB) *BanStorage {
	return &BanStorage{
		db: db,
	}
}

func (us *BanStorage) Add(ban *model.Ban) error {
	if err := us.db.Create(ban).Error; err != nil {
		return ErrBanCreate
	}
	return nil
}

func (us *BanStorage) Save(ban *model.Ban) error {
	if err := us.db.Save(ban).Error; err != nil {
		return ErrBanSave
	}
	return nil
}

func (us *BanStorage) Delete(id uint) error {
	if err := us.db.Delete(&model.Ban{}, id).Error; err != nil {
		return ErrBanDelete
	}
	return nil
}

func (us *BanStorage) Find(id uint) (*model.Ban, error) {
	ban := &model.Ban{}
	if err := us.db.First(ban, id).Error; err != nil {
		return nil, ErrBanFind
	}
	return ban, nil
}

func (us *BanStorage) FindByUserID(userID uint) ([]*model.Ban, error) {
	var bans []*model.Ban
	if err := us.db.Where("user_id = ?", userID).Find(&bans).Error; err != nil {
		return nil, ErrBanFind
	}
	return bans, nil
}

func (us *BanStorage) FindByAdminID(adminID uint) ([]*model.Ban, error) {
	var bans []*model.Ban
	if err := us.db.Where("admin_id = ?", adminID).Find(&bans).Error; err != nil {
		return nil, ErrBanFind
	}
	return bans, nil
}
