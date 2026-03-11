package upload

import (
	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/middleware"
	"github.com/parvej/luxbiss_server/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(router *gin.RouterGroup, jwtManager *jwt.Manager, rdb *redis.Client) {
	uploadGroup := router.Group("/upload")
	uploadGroup.Use(middleware.Auth(jwtManager, rdb))
	uploadGroup.Use(middleware.RequireRole("admin", "user"))

	// Create unified upload directory "public/uploads" inside project root
	// We'll pass the path accordingly
	// Since handler needs uploadDir, we instantiate here
	h := NewHandler("public/uploads")

	uploadGroup.POST("/image", h.UploadImage)
}
