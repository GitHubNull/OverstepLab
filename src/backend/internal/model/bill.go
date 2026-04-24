package model

import "time"

type Bill struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CompanyID    *uint     `json:"company_id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	Type         string    `gorm:"size:20;not null" json:"type"`
	Amount       float64   `gorm:"type:real;not null" json:"amount"`
	BalanceAfter float64   `gorm:"type:real" json:"balance_after"`
	Description  string    `gorm:"size:500" json:"description"`
	CreatedAt    time.Time `json:"created_at"`

	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (Bill) TableName() string {
	return "bills"
}
