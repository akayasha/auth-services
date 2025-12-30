package models

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserUUID  string    `gorm:"index" json:"userUuid"`
	TokenHash string    `gorm:"uniqueIndex" json:"-"` // Hide from JSON
	ExpiresAt time.Time `json:"expiresAt"`
	IsRevoked bool      `json:"isRevoked"`
	CreatedAt time.Time `json:"createdAt"`
}
