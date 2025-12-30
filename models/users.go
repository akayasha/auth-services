package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleUser     Role = "user"
	RoleEmployee Role = "employee"
)

type User struct {
	UUID string `gorm:"type:char(36);primaryKey" json:"uuid"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `gorm:"uniqueIndex" json:"username"`
	Email     string `gorm:"uniqueIndex" json:"email"`

	PasswordHash string `json:"-"` // Hide from JSON
	Role         Role   `json:"role"`

	OTPHash      string     `json:"-"` // Hide from JSON
	OTPExpiresAt *time.Time `json:"-"` // Hide from JSON

	IsEmailVerified bool `json:"isEmailVerified"`

	FailedLoginCount int        `json:"-"` // Hide from JSON
	FailedOTPCount   int        `json:"-"` // Hide from JSON
	LockedUntil      *time.Time `json:"-"` // Hide from JSON

	Dob       time.Time `json:"dob"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}
	return nil
}

func ValidateRole(role Role) error {
	switch role {
	case RoleAdmin, RoleUser, RoleEmployee:
		return nil
	default:
		return errors.New("invalid role")
	}
}
