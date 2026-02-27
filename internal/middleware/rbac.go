package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		roleSet[role] = struct{}{}
	}

	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "Authentication required",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"message":    "Invalid user role in context",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		if _, allowed := roleSet[role]; !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success":    false,
				"message":    "You do not have permission to access this resource",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		c.Next()
	}
}
