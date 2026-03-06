package giftcard

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
	giftcards := router.Group("/giftcards")
	{
		// Apply/Redeem giftcard route
		giftcards.POST("/apply", middleware.Auth(jwtManager, rdb), handler.Apply)
		// Verify giftcard (check code + balance before redeeming)
		giftcards.POST("/verify", middleware.Auth(jwtManager, rdb), handler.Verify)

		// Admin only routes
		giftcards.Use(middleware.Auth(jwtManager, rdb), middleware.RequireRole("admin"))
		{
			giftcards.POST("", handler.Create)
			giftcards.GET("", handler.List)
			giftcards.DELETE("/:id", handler.Delete)
		}
	}
}
