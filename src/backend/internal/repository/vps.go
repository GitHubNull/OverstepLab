package repository

import (
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
)

type VPSRepository struct {
	db *gorm.DB
}

var vpsRepo *VPSRepository

func GetVPSRepo() *VPSRepository {
	if vpsRepo == nil {
		vpsRepo = &VPSRepository{db: database.GetDB()}
	}
	return vpsRepo
}

func (r *VPSRepository) FindByID(id uint) (*model.VPSInstance, error) {
	var vps model.VPSInstance
	err := r.db.First(&vps, id).Error
	return &vps, err
}

func (r *VPSRepository) FindByOwnerID(ownerID uint) ([]model.VPSInstance, error) {
	var vpsList []model.VPSInstance
	err := r.db.Where("owner_id = ?", ownerID).Find(&vpsList).Error
	return vpsList, err
}

func (r *VPSRepository) FindByCompanyID(companyID uint) ([]model.VPSInstance, error) {
	var vpsList []model.VPSInstance
	err := r.db.Where("company_id = ?", companyID).Find(&vpsList).Error
	return vpsList, err
}

func (r *VPSRepository) Create(vps *model.VPSInstance) error {
	return r.db.Create(vps).Error
}

func (r *VPSRepository) Update(vps *model.VPSInstance) error {
	return r.db.Save(vps).Error
}

func (r *VPSRepository) Delete(id uint) error {
	return r.db.Delete(&model.VPSInstance{}, id).Error
}

func (r *VPSRepository) List() ([]model.VPSInstance, error) {
	var vpsList []model.VPSInstance
	err := r.db.Find(&vpsList).Error
	return vpsList, err
}
