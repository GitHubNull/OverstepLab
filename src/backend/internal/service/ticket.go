package service

import (
	"github.com/oversteplab/oversteplab/internal/model"
	"github.com/oversteplab/oversteplab/internal/repository"
	"github.com/oversteplab/oversteplab/internal/vuln"
)

type TicketService struct {
	ticketRepo    *repository.TicketRepository
	ticketReplyRepo *repository.TicketReplyRepository
}

func NewTicketService() *TicketService {
	return &TicketService{
		ticketRepo:    repository.GetTicketRepo(),
		ticketReplyRepo: repository.GetTicketReplyRepo(),
	}
}

func (s *TicketService) List(user *model.User) ([]model.Ticket, error) {
	return s.ticketRepo.FindByUserID(user.ID)
}

func (s *TicketService) Create(user *model.User, input *CreateTicketInput) (*model.Ticket, error) {
	ticket := &model.Ticket{
		Title:     input.Title,
		Content:   input.Content,
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		Status:    "open",
	}
	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

type CreateTicketInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// H-04 Vulnerability: In vulnerable mode, no ownership check
func (s *TicketService) GetDetail(user *model.User, ticketID uint) (*model.Ticket, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, err
	}

	if vuln.IsSecureMode() {
		if !user.IsPlatformAdmin() && ticket.UserID != user.ID {
			return nil, ErrUnauthorized
		}
	}

	return ticket, nil
}

func (s *TicketService) GetReplies(ticketID uint) ([]model.TicketReply, error) {
	return s.ticketReplyRepo.FindByTicketID(ticketID)
}

func (s *TicketService) Reply(user *model.User, ticketID uint, content string) (*model.TicketReply, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, err
	}

	if vuln.IsSecureMode() {
		if !user.IsPlatformAdmin() && ticket.UserID != user.ID {
			return nil, ErrUnauthorized
		}
	}

	reply := &model.TicketReply{
		TicketID: ticketID,
		UserID:   user.ID,
		Content:  content,
	}
	if err := s.ticketReplyRepo.Create(reply); err != nil {
		return nil, err
	}

	ticket.Status = "replied"
	s.ticketRepo.Update(ticket)

	return reply, nil
}

func (s *TicketService) Close(user *model.User, ticketID uint) error {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return err
	}

	if vuln.IsSecureMode() {
		if !user.IsPlatformAdmin() && ticket.UserID != user.ID {
			return ErrUnauthorized
		}
	}

	ticket.Status = "closed"
	return s.ticketRepo.Update(ticket)
}
