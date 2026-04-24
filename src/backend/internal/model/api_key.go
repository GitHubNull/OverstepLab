package model

import "time"

type APIKey struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	Name        string    `gorm:"size:100" json:"name"`
	KeyValue    string    `gorm:"size:255;not null" json:"-"`
	KeyPrefix   string    `gorm:"size:20" json:"key_prefix"`
	Permissions string    `gorm:"type:text" json:"permissions"`
	Status      string    `gorm:"size:20;default:active" json:"status"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpireAt    *time.Time `json:"expire_at"`
	CreatedAt   time.Time `json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (APIKey) TableName() string {
	return "api_keys"
}
