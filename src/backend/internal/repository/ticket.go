package repository

import (
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

var ticketRepo *TicketRepository

func GetTicketRepo() *TicketRepository {
	if ticketRepo == nil {
		ticketRepo = &TicketRepository{db: database.GetDB()}
	}
	return ticketRepo
}

func (r *TicketRepository) FindByID(id uint) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.First(&ticket, id).Error
	return &ticket, err
}

func (r *TicketRepository) FindByUserID(userID uint) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) Create(ticket *model.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *TicketRepository) Update(ticket *model.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *TicketRepository) List() ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.Order("created_at desc").Find(&tickets).Error
	return tickets, err
}

type TicketReplyRepository struct {
	db *gorm.DB
}

var ticketReplyRepo *TicketReplyRepository

func GetTicketReplyRepo() *TicketReplyRepository {
	if ticketReplyRepo == nil {
		ticketReplyRepo = &TicketReplyRepository{db: database.GetDB()}
	}
	return ticketReplyRepo
}

func (r *TicketReplyRepository) FindByTicketID(ticketID uint) ([]model.TicketReply, error) {
	var replies []model.TicketReply
	err := r.db.Where("ticket_id = ?", ticketID).Order("created_at asc").Find(&replies).Error
	return replies, err
}

func (r *TicketReplyRepository) Create(reply *model.TicketReply) error {
	return r.db.Create(reply).Error
}
