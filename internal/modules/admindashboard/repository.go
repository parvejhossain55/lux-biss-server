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

func (r *GormRepository) GetRecentActivity(ctx context.Context, limit int) ([]*ActivityResponse, error) {
	var activities []*ActivityResponse

	// Helper to format string values safely
	capitalize := func(s string) string {
		if len(s) == 0 {
			return s
		}
		// Special hack to format user statuses
		if s == "processing" {
			return "Pending"
		}
		if s == "completed" {
			return "Completed"
		}
		if s == "rejected" {
			return "Failed"
		}
		if s == "cancelled" {
			return "Failed"
		}
		if s == "active" {
			return "Active"
		}
		if s == "suspend" {
			return "Suspended"
		}
		if s == "ignored" {
			return "Ignored"
		}
		if s == "deposit" {
			return "Deposit"
		}
		if s == "withdraw" {
			return "Withdraw"
		}
		return s
	}

	// Fetch recent users
	var recentUsers []struct {
		Email     string
		Status    string
		Country   string
		CreatedAt time.Time
	}
	if err := r.db.WithContext(ctx).Table("users").Where("role = ?", user.RoleUser).Order("created_at desc").Limit(limit).Find(&recentUsers).Error; err == nil {
		for _, u := range recentUsers {
			activities = append(activities, &ActivityResponse{
				Action:     "Registration",
				Amount:     nil,
				Invoice:    "",
				Date:       u.CreatedAt.Format("01.02 03:04 PM"),
				UserStatus: capitalize(u.Status),
				Email:      u.Email,
				Country:    u.Country,
				Status:     "Completed",
				CreatedAt:  u.CreatedAt,
			})
		}
	}

	// Fetch recent transactions (with user info joined)
	var recentTxs []struct {
		Type      string
		Amount    float64
		ID        string
		CreatedAt time.Time
		TxStatus  string
		Email     string
		Country   string
		UsrStatus string
	}
	if err := r.db.WithContext(ctx).Table("transactions").
		Select("transactions.type, transactions.amount, transactions.id, transactions.created_at, transactions.status as tx_status, users.email, users.country, users.status as usr_status").
		Joins("left join users on users.id = transactions.user_id").
		Order("transactions.created_at desc").Limit(limit).Find(&recentTxs).Error; err == nil {
		for _, tx := range recentTxs {
			amt := tx.Amount
			actionType := capitalize(tx.Type)
			activities = append(activities, &ActivityResponse{
				Action:     actionType,
				Amount:     &amt,
				Invoice:    tx.ID,
				Date:       tx.CreatedAt.Format("01.02 03:04 PM"),
				UserStatus: capitalize(tx.UsrStatus),
				Email:      tx.Email,
				Country:    tx.Country,
				Status:     capitalize(tx.TxStatus),
				CreatedAt:  tx.CreatedAt,
			})
		}
	}

	// Fetch recent giftcards (with user info joined)
	var recentGiftCards []struct {
		Amount    float64
		Code      string
		CreatedAt time.Time
		TxStatus  string
		Email     string
		Country   string
		UsrStatus string
	}
	if err := r.db.WithContext(ctx).Table("giftcards").
		Select("giftcards.amount, giftcards.code, giftcards.created_at, giftcards.status as tx_status, users.email, users.country, users.status as usr_status").
		Joins("left join users on users.id = giftcards.user_id").
		Order("giftcards.created_at desc").Limit(limit).Find(&recentGiftCards).Error; err == nil {
		for _, gc := range recentGiftCards {
			amt := gc.Amount
			activities = append(activities, &ActivityResponse{
				Action:     "Gift Card",
				Amount:     &amt,
				Invoice:    gc.Code,
				Date:       gc.CreatedAt.Format("01.02 03:04 PM"),
				UserStatus: capitalize(gc.UsrStatus),
				Email:      gc.Email,
				Country:    gc.Country,
				Status:     capitalize(gc.TxStatus),
				CreatedAt:  gc.CreatedAt,
			})
		}
	}

	// Sort manually by CreatedAt desc
	for i := 0; i < len(activities)-1; i++ {
		for j := i + 1; j < len(activities); j++ {
			if activities[j].CreatedAt.After(activities[i].CreatedAt) {
				temp := activities[i]
				activities[i] = activities[j]
				activities[j] = temp
			}
		}
	}

	// Return up to limit
	if len(activities) > limit {
		activities = activities[:limit]
	}

	return activities, nil
}

var _ Repository = (*GormRepository)(nil)
