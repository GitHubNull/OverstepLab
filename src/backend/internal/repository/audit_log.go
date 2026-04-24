package repository

import (
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

var auditLogRepo *AuditLogRepository

func GetAuditLogRepo() *AuditLogRepository {
	if auditLogRepo == nil {
		auditLogRepo = &AuditLogRepository{db: database.GetDB()}
	}
	return auditLogRepo
}

func (r *AuditLogRepository) FindByUserID(userID uint) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Limit(100).Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) FindByCompanyID(companyID uint) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := r.db.Where("company_id = ?", companyID).Order("created_at desc").Limit(100).Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) List() ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := r.db.Order("created_at desc").Limit(200).Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}
