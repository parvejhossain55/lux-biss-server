package admindashboard

import (
	"context"
	"time"

	"github.com/parvej/luxbiss_server/internal/modules/transaction"
	"github.com/parvej/luxbiss_server/internal/modules/user"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetStats(ctx context.Context) (*StatsResponse, error) {
	stats := &StatsResponse{}
	todayStart := time.Now().UTC().Truncate(24 * time.Hour)
	tomorrowStart := todayStart.Add(24 * time.Hour)

	if err := r.db.WithContext(ctx).Table("users").Where("role = ?", user.RoleUser).Count(&stats.Users.Total).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("users").Where("role = ? AND created_at >= ? AND created_at < ?", user.RoleUser, todayStart, tomorrowStart).Count(&stats.Users.TodayCount).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Table("users").Where("role = ? AND status = ?", user.RoleUser, user.StatusIgnored).Count(&stats.IgnoredUsers.Total).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("users").Where("role = ? AND status = ? AND created_at >= ? AND created_at < ?", user.RoleUser, user.StatusIgnored, todayStart, tomorrowStart).Count(&stats.IgnoredUsers.TodayCount).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Table("transactions").Where("type = ?", transaction.TypeDeposit).Count(&stats.Deposits.Total).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("transactions").Where("type = ? AND created_at >= ? AND created_at < ?", transaction.TypeDeposit, todayStart, tomorrowStart).Count(&stats.Deposits.TodayCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("transactions").Where("type = ? AND created_at >= ? AND created_at < ?", transaction.TypeDeposit, todayStart, tomorrowStart).Select("COALESCE(SUM(amount), 0)").Scan(&stats.Deposits.TodayAmount).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Table("transactions").
		Where("type = ? AND status NOT IN ?", transaction.TypeWithdrawal, []string{transaction.StatusRejected, transaction.StatusCancelled}).
		Count(&stats.Withdrawals.Total).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("transactions").
		Where("type = ? AND status NOT IN ? AND created_at >= ? AND created_at < ?", transaction.TypeWithdrawal, []string{transaction.StatusRejected, transaction.StatusCancelled}, todayStart, tomorrowStart).
		Count(&stats.Withdrawals.TodayCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("transactions").
		Where("type = ? AND status NOT IN ? AND created_at >= ? AND created_at < ?", transaction.TypeWithdrawal, []string{transaction.StatusRejected, transaction.StatusCancelled}, todayStart, tomorrowStart).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.Withdrawals.TodayAmount).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Table("giftcards").Count(&stats.GiftCards.Total).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("giftcards").Where("created_at >= ? AND created_at < ?", todayStart, tomorrowStart).Count(&stats.GiftCards.TodayCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("giftcards").Where("created_at >= ? AND created_at < ?", todayStart, tomorrowStart).Select("COALESCE(SUM(amount), 0)").Scan(&stats.GiftCards.TodayAmount).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

var _ Repository = (*GormRepository)(nil)
