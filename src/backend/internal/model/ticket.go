package model

import "time"

type Ticket struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	CompanyID *uint     `json:"company_id"`
	Status    string    `gorm:"size:20;default:open" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (Ticket) TableName() string {
	return "tickets"
}

type TicketReply struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `gorm:"not null" json:"ticket_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`

	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Ticket *Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
}

func (TicketReply) TableName() string {
	return "ticket_replies"
}
