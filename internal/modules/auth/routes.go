package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/middleware"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, rdb *redis.Client) {
	auth := rg.Group("/auth")

	// Strict auth rate limit
	auth.Use(middleware.RateLimit(rdb, 10, 1*time.Minute))
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshToken)
		auth.POST("/logout", handler.Logout)
		auth.POST("/google", handler.GoogleLogin)
		auth.POST("/forgot-password", handler.ForgotPassword)
		auth.POST("/verify-otp", handler.VerifyOTP)
		auth.POST("/reset-password", handler.ResetPassword)
	}
}
