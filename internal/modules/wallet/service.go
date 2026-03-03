package wallet

import (
	"context"

	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/logger"
)

type WalletService struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) *WalletService {
	return &WalletService{repo: repo, log: log}
}

func (s *WalletService) Create(ctx context.Context, req *CreateWalletRequest) (*Wallet, error) {
	wallet := &Wallet{
		CoinName:      req.CoinName,
		Network:       req.Network,
		WalletAddress: req.WalletAddress,
		QrCodeURL:     req.QrCodeURL,
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		s.log.Errorw("Failed to create wallet", "error", err, "coin", req.CoinName)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("Wallet created successfully", "wallet_id", wallet.ID)
	return wallet, nil
}

func (s *WalletService) List(ctx context.Context) ([]*Wallet, error) {
	return s.repo.List(ctx)
}

func (s *WalletService) Update(ctx context.Context, id string, req *UpdateWalletRequest) (*Wallet, error) {
	wallet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.CoinName != nil {
		wallet.CoinName = *req.CoinName
	}
	if req.Network != nil {
		wallet.Network = *req.Network
	}
	if req.WalletAddress != nil {
		wallet.WalletAddress = *req.WalletAddress
	}
	if req.QrCodeURL != nil {
		wallet.QrCodeURL = *req.QrCodeURL
	}

	if err := s.repo.Update(ctx, wallet); err != nil {
		s.log.Errorw("Failed to update wallet", "error", err, "wallet_id", id)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("Wallet updated successfully", "wallet_id", id)
	return wallet, nil
}
