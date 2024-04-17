package models

import (
	"gorm.io/gorm"
	"neiro-api/internal/utils"
	"time"
)

type User struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Email       string `json:"email" gorm:"uniqueIndex"`
	Username    string `json:"username" gorm:"uniqueIndex"`
	DisplayName string `json:"display_name"`

	Language string `json:"language" gorm:"default:'en'"`

	Password string `json:"-"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
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
