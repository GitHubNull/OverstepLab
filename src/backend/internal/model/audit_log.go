package model

import "time"

type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	CompanyID    *uint     `json:"company_id"`
	Action       string    `gorm:"size:100;not null" json:"action"`
	ResourceType string    `gorm:"size:50" json:"resource_type"`
	ResourceID   uint      `json:"resource_id"`
	Detail       string    `gorm:"type:text" json:"detail"`
	IPAddress    string    `gorm:"size:50" json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`

	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
