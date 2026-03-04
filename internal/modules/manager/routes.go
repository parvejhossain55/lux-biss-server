package manager

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
	managers := router.Group("/managers")
	managers.Use(middleware.Auth(jwtManager, rdb), middleware.RequireRole("admin"))
	{
		managers.POST("", handler.Create)
		managers.GET("", handler.List)
		managers.GET("/:id", handler.GetByID)
		managers.PUT("/:id", handler.Update)
		managers.DELETE("/:id", handler.Delete)
	}
}
