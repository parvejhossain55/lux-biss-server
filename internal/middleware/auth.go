package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/pkg/jwt"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
)

func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "Authorization header is required",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		if !strings.HasPrefix(header, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "Authorization header must start with Bearer",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		tokenStr := strings.TrimPrefix(header, bearerPrefix)

		claims, err := jwtManager.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "Invalid or expired token",
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
