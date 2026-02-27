package health

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	health := rg.Group("/health")
	{
		health.GET("", handler.HealthCheck)
		health.GET("/ready", handler.ReadinessCheck)
	}
}
