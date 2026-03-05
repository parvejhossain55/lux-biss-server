package transaction

import (
	"context"

	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/logger"
	"github.com/parvej/luxbiss_server/internal/modules/user"
)

type TransactionService struct {
	repo        Repository
	userService user.Service
	log         *logger.Logger
}

func NewService(repo Repository, userService user.Service, log *logger.Logger) *TransactionService {
	return &TransactionService{repo: repo, userService: userService, log: log}
}

func (s *TransactionService) Create(ctx context.Context, req *CreateTransactionRequest, requestingUserID, requestingRole string) (*Transaction, error) {
	targetUserID := req.UserID

	if req.UserID == "" {
		targetUserID = requestingUserID
	}

	if targetUserID != requestingUserID && requestingRole != "admin" {
		return nil, common.ErrForbidden("You can only create transactions for your own account")
	}

	tx := &Transaction{
		UserID: targetUserID,
		Type:   req.Type,
		Amount: req.Amount,
		Status: StatusPending,
		TxHash: common.GenerateHash(),
		Note:   req.Note,
	}

	if tx.Type == TypeWithdrawal {
		u, err := s.userService.GetByID(ctx, targetUserID)
		if err != nil {
			return nil, err
		}
		if u.Balance < tx.Amount {
			return nil, common.ErrBadRequest("Insufficient balance")
		}
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		s.log.Errorw("Failed to create transaction", "error", err, "user_id", targetUserID)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("Transaction created successfully", "tx_id", tx.ID, "user_id", targetUserID)
	return tx, nil
}

func (s *TransactionService) GetByID(ctx context.Context, id string, requestingUserID, requestingRole string) (*Transaction, error) {
	tx, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if tx.UserID != requestingUserID && requestingRole != "admin" {
		return nil, common.ErrForbidden("You do not have permission to view this transaction")
	}

	return tx, nil
}

func (s *TransactionService) List(ctx context.Context, userID string, txType string, status string, limit, offset int, sortBy, sortOrder string, requestingUserID, requestingRole string) ([]*Transaction, int64, error) {
	if requestingRole != "admin" {
		userID = requestingUserID // Force user filter if not admin
	}
	return s.repo.List(ctx, userID, txType, status, limit, offset, sortBy, sortOrder)
}

func (s *TransactionService) Update(ctx context.Context, id string, req *UpdateTransactionRequest, requestingRole string) (*Transaction, error) {
	if requestingRole != "admin" {
		return nil, common.ErrForbidden("Only admins can modify transaction data")
	}
	tx, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Type != "" {
		tx.Type = req.Type
	}
	if req.Amount > 0 {
		tx.Amount = req.Amount
	}

	oldStatus := tx.Status
	newStatus := req.Status

	if newStatus != "" && newStatus != oldStatus {
		tx.Status = newStatus

		// If marked as completed, update user balance
		if newStatus == StatusCompleted {
			amountChange := tx.Amount
			if tx.Type == TypeWithdrawal {
				amountChange = -tx.Amount
			}

			if err := s.userService.UpdateBalance(ctx, tx.UserID, amountChange); err != nil {
				return nil, err
			}
		}
	}

	if req.TxHash != "" {
		tx.TxHash = req.TxHash
	}
	if req.Note != "" {
		tx.Note = req.Note
	}

	if err := s.repo.Update(ctx, tx); err != nil {
		s.log.Errorw("Failed to update transaction", "error", err, "tx_id", id)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("Transaction updated successfully", "tx_id", id)
	return tx, nil
}

func (s *TransactionService) Delete(ctx context.Context, id string, requestingRole string) error {
	if requestingRole != "admin" {
		return common.ErrForbidden("Only admins can delete transactions")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.log.Infow("Transaction deleted successfully", "tx_id", id)
	return nil
}

func (s *TransactionService) GetSummary(ctx context.Context, userID string, days int) (*Summary, error) {
	return s.repo.GetSummary(ctx, userID, days)
}
