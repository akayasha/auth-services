package services

import (
	"auth-services/config"
	"context"
	"time"
)

func BlacklistRefreshToken(tokenHash string, expiresAt time.Time) error {
	ctx := context.Background()

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}

	return config.Redis.Set(
		ctx,
		"refresh:blacklist:"+tokenHash,
		1,
		ttl,
	).Err()
}

func IsRefreshTokenBlacklisted(tokenHash string) (bool, error) {
	ctx := context.Background()

	exists, err := config.Redis.Exists(
		ctx,
		"refresh:blacklist:"+tokenHash,
	).Result()

	return exists == 1, err
}
