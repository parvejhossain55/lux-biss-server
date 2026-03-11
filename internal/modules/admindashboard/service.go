package admindashboard

import (
	"context"

	"github.com/parvej/luxbiss_server/internal/logger"
)

type DashboardService struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) *DashboardService {
	return &DashboardService{repo: repo, log: log}
}

func (s *DashboardService) GetStats(ctx context.Context) (*StatsResponse, error) {
	return s.repo.GetStats(ctx)
}

var _ Service = (*DashboardService)(nil)
