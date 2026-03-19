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

func (r *GormRepository) List(ctx context.Context, userID string, txType string, status string, limit, offset int, sortBy, sortOrder string, excludeType string) ([]*Transaction, int64, error) {
	query := r.db.WithContext(ctx).Model(&Transaction{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if txType != "" {
		query = query.Where("type = ?", txType)
	} else if excludeType != "" {
		query = query.Where("type != ?", excludeType)
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

	// Base query condition for the user (only completed transactions count towards real balance)
	baseCond := "user_id = ? AND status = ?"
	baseArgs := []interface{}{userID, StatusCompleted}

	var totalDep, totalWith, totalProf float64
	r.db.WithContext(ctx).Model(&Transaction{}).Where(baseCond+" AND type = ?", append(baseArgs, TypeDeposit)...).Select("COALESCE(SUM(amount), 0)").Scan(&totalDep)
	r.db.WithContext(ctx).Model(&Transaction{}).Where(baseCond+" AND type = ?", append(baseArgs, TypeWithdrawal)...).Select("COALESCE(SUM(amount), 0)").Scan(&totalWith)
	r.db.WithContext(ctx).Model(&Transaction{}).Where(baseCond+" AND type = ?", append(baseArgs, TypeInvestment)...).Select("COALESCE(SUM(profit_amount), 0)").Scan(&totalProf)

	var balance float64
	r.db.WithContext(ctx).Table("users").Where("id = ?", userID).Select("balance").Scan(&balance)
	availableBalance := balance

	// Calculate Withdrawal Report (Daily for the requested period)
	startDate := now.AddDate(0, 0, -days)
	var reportItems []ReportItem

	// Fetch daily withdrawal totals from DB
	var dbResults []struct {
		Date  time.Time
		Total float64
	}

	r.db.WithContext(ctx).Model(&Transaction{}).
		Select("DATE(created_at) as date, SUM(amount) as total").
		Where(baseCond+" AND type = ? AND created_at >= ?", append(baseArgs, TypeWithdrawal, startDate)...).
		Group("date").
		Order("date ASC").
		Scan(&dbResults)

	// Create a map for easy lookup
	dbMap := make(map[string]float64)
	for _, res := range dbResults {
		dbMap[res.Date.Format("2006-01-02")] = res.Total
	}

	// Fill in all days in the period (including zeros)
	for i := days; i >= 0; i-- {
		d := now.AddDate(0, 0, -i)
		dateStr := d.Format("2006-01-02")

		reportItems = append(reportItems, ReportItem{
			Date:   d.Format(time.RFC3339), // ISO 8601 format
			Amount: dbMap[dateStr],         // Will be 0 if not found in map
		})
	}

	var holdBalance float64
	r.db.WithContext(ctx).Table("users").Where("id = ?", userID).Select("hold_balance").Scan(&holdBalance)

	var withdrawableBalance float64
	r.db.WithContext(ctx).Table("users").Where("id = ?", userID).Select("withdrawable_balance").Scan(&withdrawableBalance)

	return &Summary{
		AvailableBalance:    availableBalance,
		HoldBalance:         holdBalance,
		WithdrawableBalance: withdrawableBalance,
		TotalDeposit:        totalDep,
		TotalWithdrawal:     totalWith,
		TotalProfit:         totalProf,
		PeriodDays:          days,
		WithdrawReport:      reportItems,
	}, nil
}
