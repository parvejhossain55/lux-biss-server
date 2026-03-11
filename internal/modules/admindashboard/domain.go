package admindashboard

import "context"

type Metric struct {
	Total       int64   `json:"total"`
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

type Repository interface {
	GetStats(ctx context.Context) (*StatsResponse, error)
}

type Service interface {
	GetStats(ctx context.Context) (*StatsResponse, error)
}
