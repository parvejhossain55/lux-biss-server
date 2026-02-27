package database

import (
	"context"
	"fmt"
	"time"

	"github.com/parvej/luxbiss_server/internal/config"
	"github.com/parvej/luxbiss_server/internal/logger"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.RedisConfig, log *logger.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %w", err)
	}

	log.Info("Connected to Redis",
		"host", cfg.Host,
		"port", cfg.Port,
	)

	return client, nil
}
