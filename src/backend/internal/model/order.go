package model

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderNo   string    `gorm:"size:100;uniqueIndex;not null" json:"order_no"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	CompanyID *uint     `json:"company_id"`
	VPSID     *uint     `json:"vps_id"`
	Type      string    `gorm:"size:20;not null" json:"type"`
	Amount    float64   `gorm:"type:real;not null" json:"amount"`
	Status    string    `gorm:"size:20;default:pending" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User    *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company *Company     `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	VPS     *VPSInstance `gorm:"foreignKey:VPSID" json:"vps,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}
