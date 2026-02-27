package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/config"
)

func CORS(corsConfig *config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if the origin is allowed
		if isOriginAllowed(origin, corsConfig.AllowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if len(corsConfig.AllowedOrigins) == 0 {
			// If no specific origins are set and in development, allow localhost
			if gin.Mode() == gin.DebugMode && strings.HasPrefix(origin, "http://localhost:") {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				// Reject the request
				c.AbortWithStatusJSON(403, gin.H{
					"success":    false,
					"message":    "CORS policy violation",
					"request_id": c.GetString("request_id"),
				})
				return
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", int(12*time.Hour/time.Second)))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range allowedOrigins {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}
	return false
}
