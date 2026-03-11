package transaction

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
	transactions := router.Group("/transactions")
	transactions.Use(middleware.Auth(jwtManager, rdb))
	{
		// Common endpoints (user & admin)
		transactions.POST("", handler.Create)
		transactions.GET("", handler.List)
		transactions.GET("/summary", handler.GetSummary)
		transactions.GET("/:id", handler.GetByID)
		transactions.POST("/invest", handler.Invest)

		// Admin only: modify status and delete
		adminTxs := transactions.Group("")
		adminTxs.Use(middleware.RequireRole("admin"))
		{
			adminTxs.PUT("/:id", handler.Update)
			adminTxs.DELETE("/:id", handler.Delete)
		}
	}
}
