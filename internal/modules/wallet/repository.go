package wallet

import (
	"context"
	"errors"

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

func (r *GormRepository) Create(ctx context.Context, wallet *Wallet) error {
	wallet.ID = uuid.New().String()

	result := r.db.WithContext(ctx).Create(wallet)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id string) (*Wallet, error) {
	var wallet Wallet
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&wallet)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("Wallet")
		}
		return nil, result.Error
	}

	return &wallet, nil
}

func (r *GormRepository) List(ctx context.Context) ([]*Wallet, error) {
	var wallets []*Wallet
	result := r.db.WithContext(ctx).
		Order("created_at ASC").
		Find(&wallets)
	if result.Error != nil {
		return nil, result.Error
	}

	return wallets, nil
}

func (r *GormRepository) Update(ctx context.Context, wallet *Wallet) error {
	result := r.db.WithContext(ctx).
		Model(wallet).
		Updates(map[string]interface{}{
			"coin_name":      wallet.CoinName,
			"network":        wallet.Network,
			"coin_logo_url":  wallet.CoinLogoURL,
			"wallet_address": wallet.WalletAddress,
			"qr_code_url":    wallet.QrCodeURL,
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrNotFound("Wallet")
	}

	return nil
}
