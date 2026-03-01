package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/config"
	"github.com/parvej/luxbiss_server/internal/database"
	"github.com/parvej/luxbiss_server/internal/logger"
	"github.com/parvej/luxbiss_server/internal/middleware"
	"github.com/parvej/luxbiss_server/internal/modules/auth"
	"github.com/parvej/luxbiss_server/internal/modules/health"
	"github.com/parvej/luxbiss_server/internal/modules/user"
	"github.com/parvej/luxbiss_server/internal/server"
	"github.com/parvej/luxbiss_server/pkg/email"
	"github.com/parvej/luxbiss_server/pkg/jwt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func init() {
	// Register custom validators
	common.RegisterCustomValidators(auth.RegisterPasswordValidators)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger := logger.New(cfg.Log.Level, cfg.Log.Format)
	defer appLogger.Sync()

	db, err := database.NewPostgres(&cfg.Database, appLogger)
	if err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	if err := db.AutoMigrate(&user.User{}); err != nil {
		appLogger.Fatalf("Failed to auto-migrate: %v", err)
	}
	appLogger.Info("Database migration completed")

	rdb, err := database.NewRedis(&cfg.Redis, appLogger)
	if err != nil {
		appLogger.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	jwtManager := jwt.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		cfg.JWT.Issuer,
	)

	srv := server.New(&cfg.Server, appLogger)
	router := srv.Router()

	router.Use(
		middleware.RequestID(),
		middleware.Recovery(appLogger),
		middleware.CORS(&cfg.CORS),
		middleware.CSRF(&cfg.CORS),
		middleware.SecurityHeaders(),
		middleware.RateLimit(rdb, 100, 1*time.Minute), // Global limit
	)

	api := router.Group("/api/v1")
	registerRoutes(api, db, rdb, jwtManager, appLogger, cfg)

	appLogger.Infow("Application starting",
		"app", cfg.App.Name,
		"env", cfg.App.Env,
		"port", cfg.Server.Port,
	)

	if err := srv.Start(); err != nil {
		appLogger.Fatalf("Server error: %v", err)
	}
}

func registerRoutes(
	api *gin.RouterGroup,
	db *gorm.DB,
	rdb *redis.Client,
	jwtManager *jwt.Manager,
	appLogger *logger.Logger,
	cfg *config.Config,
) {
	healthHandler := health.NewHandler()
	health.RegisterRoutes(api, healthHandler)

	userRepo := user.NewGormRepository(db)
	userService := user.NewService(userRepo, appLogger)
	userHandler := user.NewHandler(userService, appLogger)
	user.RegisterRoutes(api, userHandler, jwtManager, rdb)

	emailSender := email.NewSMTPSender(&email.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	})

	authService := auth.NewService(userService, jwtManager, rdb, emailSender, &cfg.OAuth, appLogger)
	cookieManager := auth.NewCookieManager(&cfg.Cookie, &cfg.JWT)
	authHandler := auth.NewHandler(authService, cookieManager, appLogger)
	auth.RegisterRoutes(api, authHandler, rdb)
}
