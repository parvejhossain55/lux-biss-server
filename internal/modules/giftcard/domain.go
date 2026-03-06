package giftcard

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	StatusAvailable = "Available"
	StatusUsed      = "Used"
)

type Giftcard struct {
	ID         string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	RedeemCode string         `json:"redeem_code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Amount     float64        `json:"amount" gorm:"type:decimal(10,2);not null"`
	Status     string         `json:"status" gorm:"type:varchar(20);not null;default:'Available'"`
	UserID     *string        `json:"user_id" gorm:"type:varchar(36);index"`
	UserEmail  string         `json:"user_email" gorm:"type:varchar(255)"`
	UsedAt     *time.Time     `json:"used_at"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Giftcard) TableName() string {
	return "giftcards"
}

type Repository interface {
	Create(ctx context.Context, giftcard *Giftcard) error
	GetByID(ctx context.Context, id string) (*Giftcard, error)
	GetByCode(ctx context.Context, code string) (*Giftcard, error)
	List(ctx context.Context, limit, offset int) ([]*Giftcard, int64, error)
	Update(ctx context.Context, giftcard *Giftcard) error
	Delete(ctx context.Context, id string) error
}

type Service interface {
	Create(ctx context.Context, req *CreateGiftcardRequest) (*Giftcard, error)
	List(ctx context.Context, limit, offset int) ([]*Giftcard, int64, error)
	Delete(ctx context.Context, id string) error
	Apply(ctx context.Context, req *ApplyGiftcardRequest, userID, userEmail string) (*Giftcard, error)
	Verify(ctx context.Context, req *VerifyGiftcardRequest) (*Giftcard, error)
}
