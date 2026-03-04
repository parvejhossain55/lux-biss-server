package main

import (
	"flag"
	"log"

	"github.com/parvej/luxbiss_server/internal/config"
	"github.com/parvej/luxbiss_server/internal/database"
	"github.com/parvej/luxbiss_server/internal/database/seeder"
	"github.com/parvej/luxbiss_server/internal/logger"
	"github.com/parvej/luxbiss_server/internal/modules/giftcard"
	"github.com/parvej/luxbiss_server/internal/modules/product"
	"github.com/parvej/luxbiss_server/internal/modules/user"
	"github.com/parvej/luxbiss_server/internal/modules/wallet"
)

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
		appLogger.Fatalf("Failed to get underlying sql.DB: %w", err)
	}
	defer sqlDB.Close()

	if err := db.AutoMigrate(&user.User{}, &product.Level{}, &product.Product{}, &wallet.Wallet{}, &giftcard.Giftcard{}); err != nil {
		appLogger.Fatalf("Failed to auto-migrate: %v", err)
	}

	truncate := flag.Bool("truncate", false, "truncate tables before seeding")
	flag.Parse()

	registry := seeder.NewRegistry()
	registry.Register(&seeder.UserSeeder{})
	registry.Register(&seeder.LevelSeeder{})
	registry.Register(&seeder.ProductSeeder{})
	registry.Register(&seeder.WalletSeeder{})

	if *truncate {
		appLogger.Info("Truncating tables...")
		tables := []string{"users", "levels", "products", "wallets", "giftcards"}
		if err := registry.TruncateAll(db, tables); err != nil {
			appLogger.Fatalf("Failed to truncate tables: %v", err)
		}
		appLogger.Info("Tables truncated successfully")
	}

	appLogger.Info("Starting database seeding...")
	if err := registry.RunAll(db); err != nil {
		appLogger.Fatalf("Failed to seed database: %v", err)
	}

	appLogger.Info("Database seeding completed successfully")
}
