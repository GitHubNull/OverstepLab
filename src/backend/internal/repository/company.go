package repository

import (
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

var companyRepo *CompanyRepository

func GetCompanyRepo() *CompanyRepository {
	if companyRepo == nil {
		companyRepo = &CompanyRepository{db: database.GetDB()}
	}
	return companyRepo
}

func (r *CompanyRepository) FindByID(id uint) (*model.Company, error) {
	var company model.Company
	err := r.db.First(&company, id).Error
	return &company, err
}

func (r *CompanyRepository) Create(company *model.Company) error {
	return r.db.Create(company).Error
}

func (r *CompanyRepository) Update(company *model.Company) error {
	return r.db.Save(company).Error
}

func (r *CompanyRepository) List() ([]model.Company, error) {
	var companies []model.Company
	err := r.db.Find(&companies).Error
	return companies, err
}
