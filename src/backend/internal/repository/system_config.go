package repository

import (
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
)

type SystemConfigRepository struct {
	db *gorm.DB
}

var systemConfigRepo *SystemConfigRepository

func GetSystemConfigRepo() *SystemConfigRepository {
	if systemConfigRepo == nil {
		systemConfigRepo = &SystemConfigRepository{db: database.GetDB()}
	}
	return systemConfigRepo
}

func (r *SystemConfigRepository) FindByKey(key string) (*model.SystemConfig, error) {
	var cfg model.SystemConfig
	err := r.db.Where("key = ?", key).First(&cfg).Error
	return &cfg, err
}

func (r *SystemConfigRepository) Upsert(key, value string) error {
	var cfg model.SystemConfig
	result := r.db.Where("key = ?", key).First(&cfg)
	if result.Error != nil {
		return r.db.Create(&model.SystemConfig{Key: key, Value: value}).Error
	}
	return r.db.Model(&cfg).Update("value", value).Error
}

func (r *SystemConfigRepository) List() ([]model.SystemConfig, error) {
	var configs []model.SystemConfig
	err := r.db.Order("key").Find(&configs).Error
	return configs, err
}
