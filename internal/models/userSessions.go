package models

import (
	"gorm.io/gorm"
	"time"
)

type UserSessions struct {
	ID uint `json:"id" gorm:"primaryKey"`

	RefreshToken string `json:"refresh_token" gorm:"uniqueIndex"`
	IpAddress    string `json:"ip_address"`

	UserID uint `json:"user_id"`
	User   User `json:"user" gorm:"foreignKey:user_id"`

	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
