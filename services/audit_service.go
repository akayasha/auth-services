package services

import (
	"auth-services/config"
	"auth-services/models"
	"auth-services/repository"
)

var auditRepo repository.AuditLogRepository
var auditQueue = make(chan *models.AuditLog, 1000)

func InitAuditWorker() {
	auditRepo = repository.NewAuditLogRepository(config.DB)

	go func() {
		for log := range auditQueue {
			_ = auditRepo.Create(log) // async + fire & forget
		}
	}()
}

// ✅ THIS FIXES "Unresolved reference LogAuthEvent"
func LogAuthEvent(userUUID *string, action string) {
	auditQueue <- &models.AuditLog{
		UserUUID: userUUID,
		Action:   action,
	}
}
