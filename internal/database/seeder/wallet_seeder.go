package seeder

import (
	"log"

	"github.com/google/uuid"
	"github.com/parvej/luxbiss_server/internal/modules/wallet"
	"gorm.io/gorm"
)

type WalletSeeder struct{}

func (s *WalletSeeder) Seed(db *gorm.DB) error {
	wallets := []wallet.Wallet{
		{
			ID:            uuid.New().String(),
			CoinName:      "Bitcoin",
			Network:       "Bitcoin",
			CoinLogoURL:   "",
			WalletAddress: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
			QrCodeURL:     "",
		},
		{
			ID:            uuid.New().String(),
			CoinName:      "USDT",
			Network:       "TRC20",
			CoinLogoURL:   "",
			WalletAddress: "TY8dk40k4fjs76sjhekd943kfs4",
			QrCodeURL:     "",
		},
		{
			ID:            uuid.New().String(),
			CoinName:      "Litecoin",
			Network:       "Litecoin",
			CoinLogoURL:   "",
			WalletAddress: "ltc1qzvcjm2l2cw4q8pzk5t4trm5q9u4q4gq9cxyhcf",
			QrCodeURL:     "",
		},
		{
			ID:            uuid.New().String(),
			CoinName:      "Ethereum",
			Network:       "ERC20",
			CoinLogoURL:   "",
			WalletAddress: "0x32Be343B94f860124dC4fEe278FDCBD38C102D88",
			QrCodeURL:     "",
		},
	}

	for _, w := range wallets {
		var existing wallet.Wallet
		if err := db.Where("coin_name = ? AND network = ?", w.CoinName, w.Network).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&w).Error; err != nil {
					return err
				}
				log.Printf("Seeded wallet: %s (%s)", w.CoinName, w.Network)
			} else {
				return err
			}
		}
	}

	return nil
}
