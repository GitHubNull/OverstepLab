package repository

import (
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
)

type APIKeyRepository struct {
	db *gorm.DB
}

var apiKeyRepo *APIKeyRepository

func GetAPIKeyRepo() *APIKeyRepository {
	if apiKeyRepo == nil {
		apiKeyRepo = &APIKeyRepository{db: database.GetDB()}
	}
	return apiKeyRepo
}

func (r *APIKeyRepository) FindByID(id uint) (*model.APIKey, error) {
	var key model.APIKey
	err := r.db.First(&key, id).Error
	return &key, err
}

func (r *APIKeyRepository) FindByUserID(userID uint) ([]model.APIKey, error) {
	var keys []model.APIKey
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&keys).Error
	return keys, err
}

func (r *APIKeyRepository) FindByValue(hashValue string) (*model.APIKey, error) {
	var key model.APIKey
	err := r.db.Where("key_value = ?", hashValue).First(&key).Error
	return &key, err
}

func (r *APIKeyRepository) Create(key *model.APIKey) error {
	return r.db.Create(key).Error
}

func (r *APIKeyRepository) Update(key *model.APIKey) error {
	return r.db.Save(key).Error
}

func (r *APIKeyRepository) Delete(id uint) error {
	return r.db.Delete(&model.APIKey{}, id).Error
}
