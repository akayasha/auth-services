package repository

import (
	"auth-services/models"

	"gorm.io/gorm"
)

type auditRepo struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditRepo{db: db}
}

func (r *auditRepo) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}
