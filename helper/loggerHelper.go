package helper

import (
	"auth-services/config"
	"auth-services/models"
	"context"
	"fmt"
	"time"
)

func LogAuthEvent(userID *string, event, ip, ua string) {
	log := &models.AuditLog{
		UserID:    userID,
		Event:     event,
		IP:        ip,
		UserAgent: ua,
	}
	_ = config.DB.Create(log).Error
}

func RateLimit(key string, limit int, ttl time.Duration) error {
	ctx := context.Background()
	count, _ := config.Redis.Incr(ctx, key).Result()

	if count == 1 {
		config.Redis.Expire(ctx, key, ttl)
	}

	if count > int64(limit) {
		return fmt.Errorf("too many requests")
	}

	return nil
}
