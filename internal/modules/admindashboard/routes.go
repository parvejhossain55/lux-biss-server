package admindashboard

import (
	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/middleware"
	"github.com/parvej/luxbiss_server/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
	jwtManager *jwt.Manager,
	rdb *redis.Client,
) {
	admin := router.Group("/admin/dashboard")
	admin.Use(middleware.Auth(jwtManager, rdb), middleware.RequireRole("admin"))
	{
		admin.GET("/stats", handler.GetStats)
		admin.GET("/recent-activity", handler.GetRecentActivity)
	}
}
