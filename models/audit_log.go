package models

import "time"

type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserUUID  *string   `gorm:"type:char(36);index" json:"userUuid"`
	Action    string    `gorm:"size:50;not null" json:"action"`
	IP        string    `gorm:"size:45" json:"ip"`
	UserAgent string    `gorm:"size:255" json:"userAgent"`
	CreatedAt time.Time `json:"createdAt"`
}
