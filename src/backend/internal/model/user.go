package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:100;uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Email        string    `gorm:"size:255" json:"email"`
	Phone        string    `gorm:"size:50" json:"phone"`
	Avatar       string    `gorm:"size:500" json:"avatar"`
	UserType     string    `gorm:"size:30;not null" json:"user_type"`
	CompanyID    *uint     `json:"company_id"`
	Role         string    `gorm:"size:30" json:"role"`
	Status       string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsPlatformAdmin() bool {
	return u.UserType == "platform_admin"
}

func (u *User) IsCompanyAdmin() bool {
	return u.Role == "admin" && u.UserType == "company"
}

func (u *User) IsIndividual() bool {
	return u.UserType == "individual"
}
