package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/oversteplab/oversteplab/internal/common"
	"github.com/oversteplab/oversteplab/internal/model"
	"github.com/oversteplab/oversteplab/internal/repository"
	"github.com/oversteplab/oversteplab/internal/vuln"
)

type APIKeyService struct {
	apiKeyRepo *repository.APIKeyRepository
}

func NewAPIKeyService() *APIKeyService {
	return &APIKeyService{apiKeyRepo: repository.GetAPIKeyRepo()}
}

func (s *APIKeyService) List(user *model.User) ([]model.APIKey, error) {
	return s.apiKeyRepo.FindByUserID(user.ID)
}

func (s *APIKeyService) Create(user *model.User, name string, permissions string) (*model.APIKey, error) {
	rawKey := common.GenerateAPIKey()
	hashedKey := hashAPIKey(rawKey)
	prefix := rawKey[:min(12, len(rawKey))]

	key := &model.APIKey{
		UserID:      user.ID,
		Name:        name,
		KeyValue:    hashedKey,
		KeyPrefix:   prefix,
		Permissions: permissions,
		Status:      "active",
	}

	if err := s.apiKeyRepo.Create(key); err != nil {
		return nil, err
	}

	// Return the raw key only once
	key.KeyValue = rawKey
	return key, nil
}

// H-05 Vulnerability: In vulnerable mode, no ownership check on delete
func (s *APIKeyService) Delete(user *model.User, keyID uint) error {
	key, err := s.apiKeyRepo.FindByID(keyID)
	if err != nil {
		return err
	}

	if vuln.IsSecureMode() {
		if key.UserID != user.ID && !user.IsPlatformAdmin() {
			return ErrUnauthorized
		}
	}

	return s.apiKeyRepo.Delete(keyID)
}

func hashAPIKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type BillService struct {
	billRepo    *repository.BillRepository
	companyRepo *repository.CompanyRepository
	userRepo    *repository.UserRepository
}

func NewBillService() *BillService {
	return &BillService{
		billRepo:    repository.GetBillRepo(),
		companyRepo: repository.GetCompanyRepo(),
		userRepo:    repository.GetUserRepo(),
	}
}

func (s *BillService) List(user *model.User) ([]model.Bill, error) {
	return s.billRepo.FindByUserID(user.ID)
}

func (s *BillService) Recharge(user *model.User, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}

	// Update balance
	if user.CompanyID != nil {
		company, err := s.companyRepo.FindByID(*user.CompanyID)
		if err != nil {
			return err
		}
		company.Balance += amount
		s.companyRepo.Update(company)
	}

	// Record bill
	bill := &model.Bill{
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		Type:      "recharge",
		Amount:    amount,
		Description: "Account recharge",
	}
	return s.billRepo.Create(bill)
}

func (s *BillService) Export(user *model.User) ([]model.Bill, error) {
	return s.billRepo.FindByUserID(user.ID)
}

type AuditLogService struct {
	auditRepo *repository.AuditLogRepository
}

func NewAuditLogService() *AuditLogService {
	return &AuditLogService{auditRepo: repository.GetAuditLogRepo()}
}

func (s *AuditLogService) GetLogs(user *model.User) ([]model.AuditLog, error) {
	if user.IsPlatformAdmin() {
		return s.auditRepo.List()
	}
	if user.CompanyID != nil {
		return s.auditRepo.FindByCompanyID(*user.CompanyID)
	}
	return s.auditRepo.FindByUserID(user.ID)
}

type AdminService struct {
	userRepo    *repository.UserRepository
	companyRepo *repository.CompanyRepository
	vpsRepo     *repository.VPSRepository
	auditRepo   *repository.AuditLogRepository
}

func NewAdminService() *AdminService {
	return &AdminService{
		userRepo:    repository.GetUserRepo(),
		companyRepo: repository.GetCompanyRepo(),
		vpsRepo:     repository.GetVPSRepo(),
		auditRepo:   repository.GetAuditLogRepo(),
	}
}

// V-04 Vulnerability: In vulnerable mode, no admin check (handled by middleware)
func (s *AdminService) ListUsers() ([]model.User, error) {
	users, _, err := s.userRepo.List(0, 1000)
	return users, err
}

func (s *AdminService) UpdateUserStatus(userID uint, status string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	user.Status = status
	return s.userRepo.Update(user)
}

func (s *AdminService) ListCompanies() ([]model.Company, error) {
	return s.companyRepo.List()
}

func (s *AdminService) ListAllVPS() ([]model.VPSInstance, error) {
	return s.vpsRepo.List()
}

func (s *AdminService) ListAllLogs() ([]model.AuditLog, error) {
	return s.auditRepo.List()
}
