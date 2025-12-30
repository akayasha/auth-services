package middlewares

import (
	"auth-services/config"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		key := fmt.Sprintf(
			"rate:%s:%s",
			c.ClientIP(),
			c.FullPath(),
		)

		count, err := config.Redis.Incr(ctx, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "rate limit error",
			})
			return
		}

		if count == 1 {
			config.Redis.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "too many requests",
			})
			return
		}

		c.Next()
	}
}
