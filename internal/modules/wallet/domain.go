package wallet

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID            string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	CoinName      string         `json:"coin_name" gorm:"type:varchar(50);not null;uniqueIndex:idx_coin_network"`
	Network       string         `json:"network" gorm:"type:varchar(50);not null;uniqueIndex:idx_coin_network"`
	CoinLogoURL   string         `json:"coin_logo_url" gorm:"type:varchar(255)"`
	WalletAddress string         `json:"wallet_address" gorm:"type:varchar(255);not null"`
	QrCodeURL     string         `json:"qr_code_url" gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Wallet) TableName() string {
	return "wallets"
}

type Repository interface {
	Create(ctx context.Context, wallet *Wallet) error
	GetByID(ctx context.Context, id string) (*Wallet, error)
	List(ctx context.Context) ([]*Wallet, error)
	Update(ctx context.Context, wallet *Wallet) error
}

type Service interface {
	Create(ctx context.Context, req *CreateWalletRequest) (*Wallet, error)
	List(ctx context.Context) ([]*Wallet, error)
	Update(ctx context.Context, id string, req *UpdateWalletRequest) (*Wallet, error)
}
