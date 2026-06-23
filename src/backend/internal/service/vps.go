package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/oversteplab/oversteplab/internal/common"
	"github.com/oversteplab/oversteplab/internal/model"
	"github.com/oversteplab/oversteplab/internal/repository"
	"github.com/oversteplab/oversteplab/internal/vuln"
)

var (
	ErrVPSNotFound  = errors.New("VPS not found")
	ErrUnauthorized = errors.New("unauthorized")
)

type VPSService struct {
	vpsRepo     *repository.VPSRepository
	userRepo    *repository.UserRepository
	companyRepo *repository.CompanyRepository
	orderRepo   *repository.OrderRepository
	billRepo    *repository.BillRepository
}

func NewVPSService() *VPSService {
	return &VPSService{
		vpsRepo:     repository.GetVPSRepo(),
		userRepo:    repository.GetUserRepo(),
		companyRepo: repository.GetCompanyRepo(),
		orderRepo:   repository.GetOrderRepo(),
		billRepo:    repository.GetBillRepo(),
	}
}

func (s *VPSService) List(user *model.User) ([]model.VPSInstance, error) {
	if user.IsPlatformAdmin() {
		return s.vpsRepo.List()
	}

	if user.IsCompanyAdmin() || user.Role == "operator" || user.Role == "viewer" {
		return s.vpsRepo.FindByCompanyID(*user.CompanyID)
	}

	if user.IsIndividual() {
		return s.vpsRepo.FindByOwnerID(user.ID)
	}

	return s.vpsRepo.FindByOwnerID(user.ID)
}

func (s *VPSService) GetDetail(user *model.User, vpsID uint) (*model.VPSInstance, error) {
	vps, err := s.vpsRepo.FindByID(vpsID)
	if err != nil {
		return nil, ErrVPSNotFound
	}

	// H-01 & C-01 Vulnerability: In vulnerable mode, skip ownership/company check
	if vuln.IsSecureMode() {
		if !s.canAccessVPS(user, vps) {
			return nil, ErrUnauthorized
		}
	}

	return vps, nil
}

func (s *VPSService) canAccessVPS(user *model.User, vps *model.VPSInstance) bool {
	if user.IsPlatformAdmin() {
		return true
	}
	if vps.OwnerID == user.ID {
		return true
	}
	if user.CompanyID != nil && vps.CompanyID != nil && *user.CompanyID == *vps.CompanyID {
		return true
	}
	return false
}

func (s *VPSService) canManageVPS(user *model.User) bool {
	if user.IsPlatformAdmin() {
		return true
	}
	if user.IsCompanyAdmin() {
		return true
	}
	if user.Role == "operator" {
		return true
	}
	if user.IsIndividual() {
		return true
	}
	return false
}

// V-01 Vulnerability: In vulnerable mode, no role check on power actions
func (s *VPSService) canControlVPS(user *model.User, vps *model.VPSInstance) bool {
	if vuln.IsSecureMode() {
		return s.canManageVPS(user) && s.canAccessVPS(user, vps)
	}
	// Vulnerable mode: only check access, no role check
	return true
}

func (s *VPSService) Create(user *model.User, input *CreateVPSInput) (*model.VPSInstance, error) {
	if vuln.IsSecureMode() {
		if !user.IsPlatformAdmin() && !user.IsCompanyAdmin() && !user.IsIndividual() {
			return nil, ErrUnauthorized
		}
	}

	vps := &model.VPSInstance{
		Name:      input.Name,
		OwnerID:   user.ID,
		CompanyID: user.CompanyID,
		CPU:       input.CPU,
		Memory:    input.Memory,
		Disk:      input.Disk,
		Bandwidth: input.Bandwidth,
		IPAddress: input.IPAddress,
		OSImage:   input.OSImage,
		Status:    "running",
		ExpireAt:  time.Now().AddDate(0, 1, 0),
	}

	// Auto-assign IP address if not provided
	if vps.IPAddress == "" {
		vps.IPAddress = generateRandomIP()
	}

	if err := s.vpsRepo.Create(vps); err != nil {
		return nil, err
	}

	// Create order
	order := &model.Order{
		OrderNo:   common.GenerateOrderNo(),
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		VPSID:     &vps.ID,
		Type:      "purchase",
		Amount:    float64(input.CPU)*50 + float64(input.Memory)/100,
		Status:    "paid",
	}
	s.orderRepo.Create(order)

	return vps, nil
}

type CreateVPSInput struct {
	Name      string `json:"name"`
	CPU       int    `json:"cpu"`
	Memory    int    `json:"memory"`
	Disk      int    `json:"disk"`
	Bandwidth int    `json:"bandwidth"`
	IPAddress string `json:"ip_address"`
	OSImage   string `json:"os_image"`
}

func (s *VPSService) StartVPS(user *model.User, vpsID uint) error {
	vps, err := s.vpsRepo.FindByID(vpsID)
	if err != nil {
		return ErrVPSNotFound
	}
	if !s.canControlVPS(user, vps) {
		return ErrUnauthorized
	}
	vps.Status = "running"
	return s.vpsRepo.Update(vps)
}

func (s *VPSService) StopVPS(user *model.User, vpsID uint) error {
	vps, err := s.vpsRepo.FindByID(vpsID)
	if err != nil {
		return ErrVPSNotFound
	}
	if !s.canControlVPS(user, vps) {
		return ErrUnauthorized
	}
	vps.Status = "stopped"
	return s.vpsRepo.Update(vps)
}

func (s *VPSService) RestartVPS(user *model.User, vpsID uint) error {
	vps, err := s.vpsRepo.FindByID(vpsID)
	if err != nil {
		return ErrVPSNotFound
	}
	if !s.canControlVPS(user, vps) {
		return ErrUnauthorized
	}
	vps.Status = "running"
	return s.vpsRepo.Update(vps)
}

// V-03 Vulnerability: In vulnerable mode, no role check on reinstall
func (s *VPSService) ReinstallVPS(user *model.User, vpsID uint, osImage string) error {
	vps, err := s.vpsRepo.FindByID(vpsID)
	if err != nil {
		return ErrVPSNotFound
	}

	if vuln.IsSecureMode() {
		if !user.IsPlatformAdmin() && !user.IsCompanyAdmin() && !user.IsIndividual() {
			return ErrUnauthorized
		}
	}

	if !s.canAccessVPS(user, vps) {
		return ErrUnauthorized
	}

	vps.OSImage = osImage
	return s.vpsRepo.Update(vps)
}

func (s *VPSService) DeleteVPS(user *model.User, vpsID uint) error {
	vps, err := s.vpsRepo.FindByID(vpsID)
	if err != nil {
		return ErrVPSNotFound
	}
	if vuln.IsSecureMode() {
		if !user.IsPlatformAdmin() && !user.IsCompanyAdmin() && !user.IsIndividual() {
			return ErrUnauthorized
		}
	}
	if !s.canAccessVPS(user, vps) {
		return ErrUnauthorized
	}
	return s.vpsRepo.Delete(vpsID)
}

func (s *VPSService) GetConsole(user *model.User, vpsID uint) (map[string]string, error) {
	vps, err := s.vpsRepo.FindByID(vpsID)
	if err != nil {
		return nil, ErrVPSNotFound
	}
	if vuln.IsSecureMode() {
		if !s.canManageVPS(user) || !s.canAccessVPS(user, vps) {
			return nil, ErrUnauthorized
		}
	}
	return map[string]string{
		"url":         "ws://localhost:8080/ws/console/" + hex.EncodeToString([]byte(vps.IPAddress)),
		"token":       "mock-console-token",
		"vps_id":      strconv.FormatUint(uint64(vps.ID), 10),
		"status":      vps.Status,
		"ip_address":  vps.IPAddress,
	}, nil
}

// Helper function
func hashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// generateRandomIP generates a random private IP address in 10.x.x.x range
func generateRandomIP() string {
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("10.%d.%d.%d", b[0], b[1], b[2])
}
