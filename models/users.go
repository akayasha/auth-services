package models

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleUser     Role = "user"
	RoleEmployee Role = "employee"
)

type User struct {
	UUID            string    `gorm:"type:char(36);primaryKey"`
	FirstName       string    `gorm:"type:varchar(255);not null"`
	LastName        string    `gorm:"type:varchar(255);not null"`
	Username        string    `gorm:"size:255;not null"`
	Email           string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash    string    `gorm:"not null"`
	Role            Role      `gorm:"size:50;default:'user'"`
	OTP             string    `gorm:"size:6"`
	IsEmailVerified bool      `gorm:"default:false"`
	Dob             time.Time `gorm:"not null"`
}

func ValidateRole(role Role) error {
	switch role {
	case RoleAdmin, RoleUser, RoleEmployee:
		return nil
	default:
		return errors.New("invalid role")
	}
}

// BeforeCreate hook to set UUID before saving to the database
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New().String()
	return
}
