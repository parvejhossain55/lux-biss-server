package transaction

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	TypeDeposit    = "deposit"
	TypeWithdrawal = "withdraw"
	TypeInvestment = "investment"

	StatusPending   = "processing"
	StatusCompleted = "completed"
	StatusRejected  = "rejected"
	StatusCancelled = "cancelled"
)

type Transaction struct {
	ID string `json:"id" gorm:"type:varchar(36);primaryKey"`

	UserID    string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Type      string         `json:"type" gorm:"type:varchar(20);not null;index"`
	Amount    float64        `json:"amount" gorm:"type:decimal(15,2);not null"`
	Status    string         `json:"status" gorm:"type:varchar(20);not null;default:'processing';index"`
	TxHash    string         `json:"tx_hash" gorm:"type:varchar(255)"`
	Note      string         `json:"note" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type ReportItem struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

type Summary struct {
	AvailableBalance    float64      `json:"available_balance"`
	HoldBalance         float64      `json:"hold_balance"`
	WithdrawableBalance float64      `json:"withdrawable_balance"`
	TotalDeposit        float64      `json:"total_deposit"`
	TotalWithdrawal     float64      `json:"total_withdrawal"`
	PeriodDays          int          `json:"period_days"`
	WithdrawReport      []ReportItem `json:"withdraw_report"`
}

type Repository interface {
	Create(ctx context.Context, tx *Transaction) error
	GetByID(ctx context.Context, id string) (*Transaction, error)
	List(ctx context.Context, userID string, txType string, status string, limit, offset int, sortBy, sortOrder string) ([]*Transaction, int64, error)
	Update(ctx context.Context, tx *Transaction) error
	Delete(ctx context.Context, id string) error
	GetSummary(ctx context.Context, userID string, days int) (*Summary, error)
}

type Service interface {
	Create(ctx context.Context, req *CreateTransactionRequest, requestingUserID, requestingRole string) (*Transaction, error)
	GetByID(ctx context.Context, id string, requestingUserID, requestingRole string) (*Transaction, error)
	List(ctx context.Context, userID string, txType string, status string, limit, offset int, sortBy, sortOrder string, requestingUserID, requestingRole string) ([]*Transaction, int64, error)
	Update(ctx context.Context, id string, req *UpdateTransactionRequest, requestingRole string) (*Transaction, error)
	Delete(ctx context.Context, id string, requestingRole string) error
	GetSummary(ctx context.Context, userID string, days int) (*Summary, error)

	// Investment
	Invest(ctx context.Context, userID string, req *InvestRequest) error
}
