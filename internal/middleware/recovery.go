package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/logger"
)

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("request_id")
				log.Errorw("Panic recovered",
					"error", err,
					"request_id", requestID,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success":    false,
					"message":    "Internal server error",
					"request_id": requestID,
				})
			}
		}()
		c.Next()
	}
}
