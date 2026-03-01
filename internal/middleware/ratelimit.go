package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client, rate int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := getRateLimitKey(c)
		rKey := fmt.Sprintf("ratelimit:%s:%s", c.Request.URL.Path, key)

		// Using Redis INCR and EXPIRE for sliding window (approx)
		count, err := rdb.Incr(c.Request.Context(), rKey).Result()
		if err != nil {
			c.Next() // Fallback: Allow if Redis is down
			return
		}

		if count == 1 {
			rdb.Expire(c.Request.Context(), rKey, window)
		}

		if count > int64(rate) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success":    false,
				"message":    "Too many requests, please try again later",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		c.Next()
	}
}

func getRateLimitKey(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return "user:" + userID.(string)
	}
	return "ip:" + c.ClientIP()
}
