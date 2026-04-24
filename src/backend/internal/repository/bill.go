package repository

import (
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
)

type BillRepository struct {
	db *gorm.DB
}

var billRepo *BillRepository

func GetBillRepo() *BillRepository {
	if billRepo == nil {
		billRepo = &BillRepository{db: database.GetDB()}
	}
	return billRepo
}

func (r *BillRepository) FindByUserID(userID uint) ([]model.Bill, error) {
	var bills []model.Bill
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&bills).Error
	return bills, err
}

func (r *BillRepository) Create(bill *model.Bill) error {
	return r.db.Create(bill).Error
}
