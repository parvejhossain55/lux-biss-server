package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int
	window   time.Duration
}

type visitor struct {
	count    int
	lastSeen time.Time
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get user ID first, fallback to IP
		key := rl.getRateLimitKey(c)

		rl.mu.Lock()
		v, exists := rl.visitors[key]
		if !exists || time.Since(v.lastSeen) > rl.window {
			rl.visitors[key] = &visitor{count: 1, lastSeen: time.Now()}
			rl.mu.Unlock()
			c.Next()
			return
		}

		v.count++
		v.lastSeen = time.Now()

		if v.count > rl.rate {
			rl.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success":    false,
				"message":    "Too many requests, please try again later",
				"request_id": c.GetString("request_id"),
			})
			return
		}
		rl.mu.Unlock()

		c.Next()
	}
}

// getRateLimitKey returns a unique key for rate limiting based on user authentication
func (rl *RateLimiter) getRateLimitKey(c *gin.Context) string {
	// If user is authenticated, use user ID for per-user rate limiting
	if userID, exists := c.Get("user_id"); exists {
		return "user:" + userID.(string)
	}

	// Fallback to IP-based rate limiting for unauthenticated requests
	return "ip:" + c.ClientIP()
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
