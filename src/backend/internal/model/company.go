package model

import "time"

type Company struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	LicenseNo    string    `gorm:"size:100" json:"license_no"`
	ContactName  string    `gorm:"size:100" json:"contact_name"`
	ContactPhone string    `gorm:"size:50" json:"contact_phone"`
	Balance      float64   `gorm:"type:real;default:0" json:"balance"`
	Status       string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Company) TableName() string {
	return "companies"
}
