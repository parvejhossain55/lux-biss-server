package wallet

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
	wallets := router.Group("/wallets")
	{
		// Public routes (if any, typically getting list might be public or for logged-in users)
		wallets.GET("", middleware.Auth(jwtManager, rdb), handler.List)

		// Admin only routes
		wallets.Use(middleware.Auth(jwtManager, rdb), middleware.RequireRole("admin"))
		{
			wallets.POST("", handler.Create)
			wallets.PUT("/:id", handler.Update)
			wallets.DELETE("/:id", handler.Delete)
		}
	}
}
