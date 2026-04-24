package repository

import (
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

var orderRepo *OrderRepository

func GetOrderRepo() *OrderRepository {
	if orderRepo == nil {
		orderRepo = &OrderRepository{db: database.GetDB()}
	}
	return orderRepo
}

func (r *OrderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.First(&order, id).Error
	return &order, err
}

func (r *OrderRepository) FindByUserID(userID uint) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}
