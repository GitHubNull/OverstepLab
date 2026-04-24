package repository

import (
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var userRepo *UserRepository

func GetUserRepo() *UserRepository {
	if userRepo == nil {
		userRepo = &UserRepository{db: database.GetDB()}
	}
	return userRepo
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByCompanyID(companyID uint) ([]model.User, error) {
	var users []model.User
	err := r.db.Where("company_id = ?", companyID).Find(&users).Error
	return users, err
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepository) List(offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	r.db.Model(&model.User{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}
