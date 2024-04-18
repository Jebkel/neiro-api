package models

import (
	"gorm.io/gorm"
	"neiro-api/internal/utils"
	"time"
)

type User struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Email           string     `json:"email" gorm:"uniqueIndex"`
	Username        string     `json:"username" gorm:"uniqueIndex"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" gorm:"type:TIMESTAMP;null;default:null"`

	Language string `json:"language" gorm:"default:'en'"`

	Password string `json:"-"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

var (
	argonP = &utils.ArgonParams{
		Memory:      64 * 1024,
		Iterations:  2,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
)

func (u *User) HashPassword() error {
	encodedHash, err := utils.GenerateFromPassword(u.Password, argonP)
	if err != nil {
		return err
	}
	u.Password = encodedHash
	return nil
}

func (u *User) ValidatePassword(password string) (bool, error) {
	match, err := utils.ComparePasswordAndHash(password, u.Password)
	if err != nil {
		return false, err
	}
	return match, nil
}
