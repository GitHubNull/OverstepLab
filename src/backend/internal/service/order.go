package service

import (
	"github.com/oversteplab/oversteplab/internal/model"
	"github.com/oversteplab/oversteplab/internal/repository"
	"github.com/oversteplab/oversteplab/internal/vuln"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
}

func NewOrderService() *OrderService {
	return &OrderService{orderRepo: repository.GetOrderRepo()}
}

func (s *OrderService) List(user *model.User) ([]model.Order, error) {
	return s.orderRepo.FindByUserID(user.ID)
}

// H-03 Vulnerability: In vulnerable mode, no ownership check
func (s *OrderService) GetDetail(user *model.User, orderID uint) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if vuln.IsSecureMode() {
		if !user.IsPlatformAdmin() && order.UserID != user.ID {
			return nil, ErrUnauthorized
		}
	}

	return order, nil
}
