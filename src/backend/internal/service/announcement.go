package service

import (
	"errors"

	"github.com/oversteplab/oversteplab/internal/model"
	"github.com/oversteplab/oversteplab/internal/repository"
)

type AnnouncementService struct {
	repo     *repository.AnnouncementRepository
	userRepo *repository.UserRepository
}

func NewAnnouncementService() *AnnouncementService {
	return &AnnouncementService{
		repo:     repository.GetAnnouncementRepo(),
		userRepo: repository.GetUserRepo(),
	}
}

func (s *AnnouncementService) Create(user *model.User, title, content string, isPinned bool) (*model.Announcement, error) {
	if !user.IsPlatformAdmin() {
		return nil, ErrUnauthorized
	}
	a := &model.Announcement{
		Title:    title,
		Content:  content,
		UserID:   user.ID,
		IsPinned: isPinned,
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AnnouncementService) Update(user *model.User, id uint, title, content string, isPinned bool) error {
	if !user.IsPlatformAdmin() {
		return ErrUnauthorized
	}
	a, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("announcement not found")
	}
	a.Title = title
	a.Content = content
	a.IsPinned = isPinned
	return s.repo.Update(a)
}

func (s *AnnouncementService) Delete(user *model.User, id uint) error {
	if !user.IsPlatformAdmin() {
		return ErrUnauthorized
	}
	return s.repo.Delete(id)
}

func (s *AnnouncementService) List() ([]model.Announcement, error) {
	return s.repo.List()
}
