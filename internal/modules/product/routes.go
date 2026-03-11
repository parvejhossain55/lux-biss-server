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

	// Levels
	rg.GET("/levels", handler.ListLevels)
	rg.GET("/levels/:level_id/steps", handler.ListStepsByLevel)

	// Admin only routes
	protected := rg.Group("")
	protected.Use(middleware.Auth(jwtManager, rdb))
	protected.Use(middleware.RequireRole("admin"))
	{
		// Products
		protected.POST("/products", handler.Create)
		protected.PUT("/products/:id", handler.Update)
		protected.DELETE("/products/:id", handler.Delete)

		// Levels
		protected.POST("/levels", handler.CreateLevel)
		protected.PUT("/levels/:id", handler.UpdateLevel)
		protected.DELETE("/levels/:id", handler.DeleteLevel)

		// Steps
		protected.POST("/steps", handler.CreateStep)
		protected.PUT("/steps/:id", handler.UpdateStep)
		protected.DELETE("/steps/:id", handler.DeleteStep)
	}
}
