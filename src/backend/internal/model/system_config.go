package model

type SystemConfig struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Key   string `json:"key" gorm:"uniqueIndex;not null"`
	Value string `json:"value" gorm:"not null;default:''"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}
