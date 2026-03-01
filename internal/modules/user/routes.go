package user

import (
	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/middleware"
	"github.com/parvej/luxbiss_server/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwt.Manager, rdb *redis.Client) {
	users := rg.Group("/users")

	users.Use(middleware.Auth(jwtManager, rdb))
	{
		users.GET("/me", handler.GetMe)
		users.GET("/:id", handler.GetByID)
		users.GET("", middleware.RequireRole("admin"), handler.List)
		users.PUT("/:id", handler.Update)
		users.DELETE("/:id", middleware.RequireRole("admin"), handler.Delete)
	}
}
