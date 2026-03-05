package transaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/parvej/luxbiss_server/internal/common"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, tx *Transaction) error {
	tx.ID = uuid.New().String()

	if tx.Status == "" {
		tx.Status = StatusPending
	}

	result := r.db.WithContext(ctx).Create(tx)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id string) (*Transaction, error) {
	var tx Transaction
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&tx)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("Transaction")
		}
		return nil, result.Error
	}

	return &tx, nil
}

func (r *GormRepository) List(ctx context.Context, userID string, txType string, status string, limit, offset int, sortBy, sortOrder string) ([]*Transaction, int64, error) {
	query := r.db.WithContext(ctx).Model(&Transaction{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Validate sortBy to allow only specific columns
	validSortColumns := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"amount":     true,
		"status":     true,
		"type":       true,
	}

	if !validSortColumns[sortBy] {
		sortBy = "created_at"
	}

	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	var txs []*Transaction
	result := query.
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&txs)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return txs, total, nil
}

func (r *GormRepository) Update(ctx context.Context, tx *Transaction) error {
	result := r.db.WithContext(ctx).Save(tx)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound("Transaction")
	}
	return nil
}

func (r *GormRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&Transaction{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound("Transaction")
	}
	return nil
}

func (r *GormRepository) GetSummary(ctx context.Context, userID string, days int) (*Summary, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -days)
	prevStartDate := startDate.AddDate(0, 0, -days)

	// Base query condition for the user (only completed transactions count towards real balance)
	baseCond := "user_id = ? AND status = ?"
	baseArgs := []interface{}{userID, StatusCompleted}

	var totalDep, totalWith float64
	r.db.WithContext(ctx).Model(&Transaction{}).Where(baseCond+" AND type = ?", append(baseArgs, TypeDeposit)...).Select("COALESCE(SUM(amount), 0)").Scan(&totalDep)
	r.db.WithContext(ctx).Model(&Transaction{}).Where(baseCond+" AND type = ?", append(baseArgs, TypeWithdrawal)...).Select("COALESCE(SUM(amount), 0)").Scan(&totalWith)
	availableBalance := totalDep - totalWith

	// Calculate Current Period Stats
	var currentDep, currentWith float64
	r.db.WithContext(ctx).Model(&Transaction{}).Where(baseCond+" AND type = ? AND created_at >= ?", append(baseArgs, TypeDeposit, startDate)...).Select("COALESCE(SUM(amount), 0)").Scan(&currentDep)
	r.db.WithContext(ctx).Model(&Transaction{}).Where(baseCond+" AND type = ? AND created_at >= ?", append(baseArgs, TypeWithdrawal, startDate)...).Select("COALESCE(SUM(amount), 0)").Scan(&currentWith)

	// Calculate Previous Period Stats
	var prevDep, prevWith float64
	r.db.WithContext(ctx).Model(&Transaction{}).Where(baseCond+" AND type = ? AND created_at >= ? AND created_at < ?", append(baseArgs, TypeDeposit, prevStartDate, startDate)...).Select("COALESCE(SUM(amount), 0)").Scan(&prevDep)
	r.db.WithContext(ctx).Model(&Transaction{}).Where(baseCond+" AND type = ? AND created_at >= ? AND created_at < ?", append(baseArgs, TypeWithdrawal, prevStartDate, startDate)...).Select("COALESCE(SUM(amount), 0)").Scan(&prevWith)

	// Percentage Change calculation
	depChange := calculatePercentageChange(currentDep, prevDep)
	withChange := calculatePercentageChange(currentWith, prevWith)

	// Calculate Withdrawal Report (Always last 6 months)
	reportStartDate := now.AddDate(0, -6, 0)
	var reportItems []ReportItem
	var dbResults []struct {
		Label string
		Total float64
	}

	r.db.WithContext(ctx).Model(&Transaction{}).
		Select("TO_CHAR(created_at, 'Mon') as label, TO_CHAR(created_at, 'YYYY-MM') as sort_key, SUM(amount) as total").
		Where(baseCond+" AND type = ? AND created_at >= ?", append(baseArgs, TypeWithdrawal, reportStartDate)...).
		Group("label, sort_key").
		Order("sort_key ASC").
		Scan(&dbResults)

	for _, res := range dbResults {
		reportItems = append(reportItems, ReportItem{
			Month: res.Label,
			Value: res.Total,
		})
	}

	return &Summary{
		AvailableBalance: availableBalance,
		TotalDeposit:     totalDep,
		TotalWithdrawal:  totalWith,
		DepositChange:    depChange,
		WithdrawalChange: withChange,
		PeriodDays:       days,
		WithdrawReport:   reportItems,
	}, nil
}

func calculatePercentageChange(current, previous float64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100.0 // 100% increase if prev was 0
		}
		return 0.0
	}
	return ((current - previous) / previous) * 100.0
}
