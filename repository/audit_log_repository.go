package repository

import "auth-services/models"

type AuditLogRepository interface {
	Create(log *models.AuditLog) error
}
