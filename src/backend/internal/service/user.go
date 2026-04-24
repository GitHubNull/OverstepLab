package service

import (
	"github.com/oversteplab/oversteplab/internal/common"
	"github.com/oversteplab/oversteplab/internal/model"
	"github.com/oversteplab/oversteplab/internal/repository"
	"github.com/oversteplab/oversteplab/internal/vuln"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{userRepo: repository.GetUserRepo()}
}

func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) UpdateProfile(userID uint, email, phone string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if email != "" {
		user.Email = email
	}
	if phone != "" {
		user.Phone = phone
	}
	return s.userRepo.Update(user)
}

func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if !common.CheckPassword(oldPassword, user.PasswordHash) {
		return ErrInvalidCredentials
	}
	hash, err := common.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.userRepo.Update(user)
}

// H-02 Vulnerability: In vulnerable mode, allow viewing other users' profiles
func (s *UserService) GetUserByID(currentUser *model.User, targetID uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(targetID)
	if err != nil {
		return nil, err
	}

	if vuln.IsSecureMode() {
		if currentUser.ID != user.ID && !currentUser.IsPlatformAdmin() {
			return nil, ErrUnauthorized
		}
	}

	return user, nil
}
