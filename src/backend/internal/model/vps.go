package model

import "time"

type VPSInstance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	OwnerID   uint      `gorm:"not null" json:"owner_id"`
	CompanyID *uint     `json:"company_id"`
	CPU       int       `json:"cpu"`
	Memory    int       `json:"memory"`
	Disk      int       `json:"disk"`
	Bandwidth int       `json:"bandwidth"`
	IPAddress string    `gorm:"size:50" json:"ip_address"`
	OSImage   string    `gorm:"size:100" json:"os_image"`
	Status    string    `gorm:"size:20;default:running" json:"status"`
	ExpireAt  time.Time `json:"expire_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Owner   *User    `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (VPSInstance) TableName() string {
	return "vps_instances"
}
