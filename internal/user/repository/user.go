package repository

import (
	"log"
	"time"

	"github.com/SocialNetworkY/Backend/internal/user/model"
	"github.com/SocialNetworkY/Backend/pkg/constant"
	"gorm.io/gorm"
)

type (
	UserSearch interface {
		Index(user *model.User) error
		Delete(id uint) error
		Search(query string, skip, limit int) ([]uint, error)
	}

	UserRepository struct {
		db *gorm.DB
		s  UserSearch
	}
)

func NewUserRepository(db *gorm.DB, s UserSearch) *UserRepository {
	return &UserRepository{
		db: db,
		s:  s,
	}
}

func (ur *UserRepository) Add(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return ur.s.Index(user)
}

func (ur *UserRepository) Save(user *model.User) error {
	if err := ur.db.Save(user).Error; err != nil {
		return err
	}
	return ur.s.Index(user)
}

func (ur *UserRepository) Delete(user *model.User) error {
	if err := ur.db.Delete(user).Error; err != nil {
		return err
	}
	return ur.s.Delete(user.ID)
}

func (ur *UserRepository) ExistsByLogin(login string) (bool, error) {
	exists, err := ur.ExistsByEmail(login)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	exists, err = ur.ExistsByUsername(login)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (ur *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := ur.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ur *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := ur.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ur *UserRepository) FindByLogin(login string) (*model.User, error) {
	user, err := ur.FindByEmail(login)
	if err == nil {
		return user, nil
	}

	user, err = ur.FindByUsername(login)
	if err == nil {
		return user, nil
	}

	return nil, err
}

func (ur *UserRepository) Find(id uint) (*model.User, error) {
	user := &model.User{}
	if err := ur.db.Preload("Bans").Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) FindSome(skip, limit int) ([]*model.User, error) {
	var users []*model.User
	if limit < 0 {
		skip = -1
	}
	if err := ur.db.Preload("Bans").Offset(skip).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *UserRepository) FindByUsername(username string) (*model.User, error) {
	user := &model.User{}
	if err := ur.db.Preload("Bans").Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := ur.db.Preload("Bans").Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) FindByNickname(nickname string, skip, limit int) ([]*model.User, error) {
	var users []*model.User
	if limit < 0 {
		skip = -1
	}
	if err := ur.db.Preload("Bans").Where("nickname = ?", nickname).Offset(skip).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *UserRepository) Search(query string, skip, limit int) ([]*model.User, error) {
	usersIDs, err := ur.s.Search(query, skip, limit)
	if err != nil {
		return nil, err
	}

	var users []*model.User
	if err := ur.db.Preload("Bans").Where("id IN ?", usersIDs).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) Statistic() (*model.UserStatistic, error) {
	log.Println("Getting user statistics")

	var stat model.UserStatistic
	err := ur.db.Model(&model.User{}).
		Select("COUNT(*) AS total, "+
			"SUM(CASE WHEN role > ? THEN 1 ELSE 0 END) AS admin, "+
			"SUM(CASE WHEN EXISTS ("+
			"  SELECT 1 FROM bans WHERE bans.user_id = users.id AND bans.expired_at > ? AND bans.unban_reason = ''"+
			") THEN 1 ELSE 0 END) AS banned, "+
			"SUM(CASE WHEN NOT EXISTS ("+
			"  SELECT 1 FROM bans WHERE bans.user_id = users.id AND bans.expired_at > ? AND bans.unban_reason = ''"+
			") THEN 1 ELSE 0 END) AS active",
			constant.RoleUser, time.Now(), time.Now()).
		Scan(&stat).Error

	if err != nil {
		log.Printf("Error getting user statistics: %v\n", err)
		return nil, err
	}

	log.Printf("User statistics found: %+v\n", stat)
	return &stat, nil
}
