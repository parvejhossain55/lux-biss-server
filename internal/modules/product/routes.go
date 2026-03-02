package product

import (
	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/middleware"
	"github.com/parvej/luxbiss_server/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwt.Manager, rdb *redis.Client) {
	products := rg.Group("/products")

	// Public routes
	products.GET("", handler.List)
	products.GET("/:id", handler.GetByID)

	// Admin only routes
	protected := products.Group("")
	protected.Use(middleware.Auth(jwtManager, rdb))
	protected.Use(middleware.RequireRole("admin"))
	{
		protected.POST("", handler.Create)
		protected.PUT("/:id", handler.Update)
		protected.DELETE("/:id", handler.Delete)
	}
}
