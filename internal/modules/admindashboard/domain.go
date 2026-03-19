package admindashboard

import (
	"context"
	"time"
)

type Metric struct {
	Total       int64   `json:"total"`
	TotalAmount float64 `json:"total_amount,omitempty"`
	TodayCount  int64   `json:"today_count"`
	TodayAmount float64 `json:"today_amount,omitempty"`
}

type StatsResponse struct {
	Users        Metric `json:"users"`
	IgnoredUsers Metric `json:"ignored_users"`
	Deposits     Metric `json:"deposits"`
	Withdrawals  Metric `json:"withdrawals"`
	GiftCards    Metric `json:"gift_cards"`
}

type ActivityResponse struct {
	Action     string    `json:"action"` // "Deposit", "Withdraw", "Gift Card", "Registration"
	Amount     *float64  `json:"amount"` // nullable for registration
	Invoice    string    `json:"invoice"`
	Date       string    `json:"date"`
	UserStatus string    `json:"userStatus"`
	Email      string    `json:"email"`
	Country    string    `json:"country"`
	Status     string    `json:"status"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"-"` // Used for sorting before responding
}

type Repository interface {
	GetStats(ctx context.Context) (*StatsResponse, error)
	GetRecentActivity(ctx context.Context, limit int) ([]*ActivityResponse, error)
}

type Service interface {
	GetStats(ctx context.Context) (*StatsResponse, error)
	GetRecentActivity(ctx context.Context, limit int) ([]*ActivityResponse, error)
}
