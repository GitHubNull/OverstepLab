package repository

import (
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
)

type AnnouncementRepository struct {
	db *gorm.DB
}

var announcementRepo *AnnouncementRepository

func GetAnnouncementRepo() *AnnouncementRepository {
	if announcementRepo == nil {
		announcementRepo = &AnnouncementRepository{db: database.GetDB()}
	}
	return announcementRepo
}

func (r *AnnouncementRepository) FindByID(id uint) (*model.Announcement, error) {
	var a model.Announcement
	err := r.db.Preload("User").First(&a, id).Error
	return &a, err
}

func (r *AnnouncementRepository) Create(a *model.Announcement) error {
	return r.db.Create(a).Error
}

func (r *AnnouncementRepository) Update(a *model.Announcement) error {
	return r.db.Save(a).Error
}

func (r *AnnouncementRepository) Delete(id uint) error {
	return r.db.Delete(&model.Announcement{}, id).Error
}

func (r *AnnouncementRepository) List() ([]model.Announcement, error) {
	var list []model.Announcement
	err := r.db.Preload("User").Order("is_pinned DESC, created_at DESC").Find(&list).Error
	return list, err
}
