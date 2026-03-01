package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
)

func Auth(jwtManager *jwt.Manager, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		// 1. Try to read access token from cookie first
		if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
			tokenStr = cookie
		}

		// 2. Fallback to Authorization header (for API clients, Postman, mobile apps, etc.)
		if tokenStr == "" {
			header := c.GetHeader(authorizationHeader)
			if header != "" {
				if !strings.HasPrefix(header, bearerPrefix) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"success":    false,
						"message":    "Authorization header must start with Bearer",
						"request_id": c.GetString("request_id"),
					})
					return
				}
				tokenStr = strings.TrimPrefix(header, bearerPrefix)
			}
		}

		// 3. No token found from either source
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "Authentication required",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		// 4. Validate the token
		claims, err := jwtManager.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "Invalid or expired token",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		// 5. Check if token is blacklisted (Logout check)
		blKey := fmt.Sprintf("blacklist:%s", claims.ID)
		blacklisted, _ := rdb.Exists(c.Request.Context(), blKey).Result()
		if blacklisted > 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "Token has been revoked",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}
