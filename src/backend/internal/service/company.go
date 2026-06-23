package service

import (
	"errors"

	"github.com/oversteplab/oversteplab/internal/common"
	"github.com/oversteplab/oversteplab/internal/model"
	"github.com/oversteplab/oversteplab/internal/repository"
	"github.com/oversteplab/oversteplab/internal/vuln"
)

var ErrMemberNotFound = errors.New("member not found")

type CompanyService struct {
	userRepo    *repository.UserRepository
	companyRepo *repository.CompanyRepository
}

func NewCompanyService() *CompanyService {
	return &CompanyService{
		userRepo:    repository.GetUserRepo(),
		companyRepo: repository.GetCompanyRepo(),
	}
}

func (s *CompanyService) ListMembers(user *model.User) ([]model.User, error) {
	if user.CompanyID == nil {
		return nil, ErrUnauthorized
	}
	return s.userRepo.FindByCompanyID(*user.CompanyID)
}

// V-02 & C-02 Vulnerability: In vulnerable mode, no role/type check
func (s *CompanyService) AddMember(user *model.User, input *AddMemberInput) (*model.User, error) {
	if vuln.IsSecureMode() {
		// Secure mode: only company admin can add members
		if !user.IsCompanyAdmin() && !user.IsPlatformAdmin() {
			return nil, ErrUnauthorized
		}
		// C-02: Individuals cannot create company members
		if user.IsIndividual() {
			return nil, ErrUnauthorized
		}
	}

	hash, err := common.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	member := &model.User{
		Username:     input.Username,
		PasswordHash: hash,
		Email:        input.Email,
		Phone:        input.Phone,
		UserType:     "company",
		CompanyID:    user.CompanyID,
		Role:         input.Role,
		Status:       "active",
	}

	if err := s.userRepo.Create(member); err != nil {
		return nil, err
	}

	return member, nil
}

type AddMemberInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

func (s *CompanyService) UpdateMember(user *model.User, memberID uint, input *UpdateMemberInput) error {
	member, err := s.userRepo.FindByID(memberID)
	if err != nil {
		return ErrMemberNotFound
	}

	if vuln.IsSecureMode() {
		if !user.IsCompanyAdmin() && !user.IsPlatformAdmin() {
			return ErrUnauthorized
		}
		if member.CompanyID == nil || user.CompanyID == nil || *member.CompanyID != *user.CompanyID {
			return ErrUnauthorized
		}
	}

	if input.Email != "" {
		member.Email = input.Email
	}
	if input.Phone != "" {
		member.Phone = input.Phone
	}
	if input.Status != "" {
		member.Status = input.Status
	}

	return s.userRepo.Update(member)
}

type UpdateMemberInput struct {
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

func (s *CompanyService) DeleteMember(user *model.User, memberID uint) error {
	member, err := s.userRepo.FindByID(memberID)
	if err != nil {
		return ErrMemberNotFound
	}

	if vuln.IsSecureMode() {
		if !user.IsCompanyAdmin() && !user.IsPlatformAdmin() {
			return ErrUnauthorized
		}
		if member.CompanyID == nil || user.CompanyID == nil || *member.CompanyID != *user.CompanyID {
			return ErrUnauthorized
		}
	}

	return s.userRepo.Delete(memberID)
}

// V-05 Vulnerability: In vulnerable mode, users can change their own role
func (s *CompanyService) ChangeRole(user *model.User, targetID uint, newRole string) error {
	target, err := s.userRepo.FindByID(targetID)
	if err != nil {
		return ErrMemberNotFound
	}

	if vuln.IsSecureMode() {
		if !user.IsCompanyAdmin() && !user.IsPlatformAdmin() {
			return ErrUnauthorized
		}
		if target.CompanyID == nil || user.CompanyID == nil || *target.CompanyID != *user.CompanyID {
			return ErrUnauthorized
		}
	}

	// In vulnerable mode, no check - user can change own role (V-05)
	target.Role = newRole
	return s.userRepo.Update(target)
}
