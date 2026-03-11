package transaction

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/logger"
	"github.com/parvej/luxbiss_server/internal/modules/product"
	"github.com/parvej/luxbiss_server/internal/modules/user"
)

type TransactionService struct {
	repo           Repository
	productService product.Service
	userService    user.Service
	log            *logger.Logger
}

func NewService(repo Repository, userService user.Service, productService product.Service, log *logger.Logger) *TransactionService {
	return &TransactionService{repo: repo, userService: userService, productService: productService, log: log}
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
		if u.WithdrawalAddress == "" {
			return nil, common.ErrBadRequest("Please set your withdrawal address in your profile before requesting a withdrawal")
		}
		if u.WithdrawableBalance < tx.Amount {
			return nil, common.ErrBadRequest("Insufficient balance")
		}
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		s.log.Errorw("Failed to create transaction", "error", err, "user_id", targetUserID)
		return nil, common.ErrInternal(err)
	}

	// For deposits, add to hold balance until approved
	if tx.Type == TypeDeposit {
		if err := s.userService.UpdateHoldBalance(ctx, targetUserID, tx.Amount); err != nil {
			s.log.Errorw("Failed to update user hold balance on deposit creation", "error", err, "user_id", targetUserID)
		}
	}
	// For withdrawals, reserve withdrawable balance immediately
	if tx.Type == TypeWithdrawal {
		if err := s.userService.UpdateWithdrawableBalance(ctx, targetUserID, -tx.Amount); err != nil {
			s.log.Errorw("Failed to reserve withdrawable balance on withdrawal creation", "error", err, "user_id", targetUserID, "tx_id", tx.ID)
			_ = s.repo.Delete(ctx, tx.ID)
			return nil, common.ErrInternal(err)
		}
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
		if newStatus == StatusCompleted && oldStatus == StatusPending {
			if tx.Type == TypeDeposit {
				// Move from HoldBalance to main Balance
				if err := s.userService.UpdateHoldBalance(ctx, tx.UserID, -tx.Amount); err != nil {
					s.log.Errorw("Failed to update user hold balance", "error", err, "tx_id", id)
				}
				if err := s.userService.UpdateBalance(ctx, tx.UserID, tx.Amount); err != nil {
					return nil, err
				}
			}
		} else if (newStatus == StatusRejected || newStatus == StatusCancelled) && oldStatus == StatusPending {
			if tx.Type == TypeDeposit {
				// Remove from HoldBalance
				if err := s.userService.UpdateHoldBalance(ctx, tx.UserID, -tx.Amount); err != nil {
					s.log.Errorw("Failed to update user hold balance on rejection", "error", err, "tx_id", id)
				}
			} else if tx.Type == TypeWithdrawal {
				// Release reserved withdrawable balance
				if err := s.userService.UpdateWithdrawableBalance(ctx, tx.UserID, tx.Amount); err != nil {
					s.log.Errorw("Failed to release withdrawable balance on rejection", "error", err, "tx_id", id)
				}
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

func (s *TransactionService) Invest(ctx context.Context, userID string, req *InvestRequest) error {
	// 1. Get step products
	products, _, err := s.productService.List(ctx, 100, 0, "", "", req.LevelID, req.StepID)
	if err != nil {
		return err
	}
	if len(products) == 0 {
		return common.ErrNotFound("No products found for this step")
	}

	// 2. Calculate total cost
	totalCost := 0.0
	for _, p := range products {
		minQty := p.MinQuantity
		if minQty <= 0 {
			minQty = 1
		}
		totalCost += p.Price * float64(minQty)
	}
	totalCost *= float64(req.Quantity)

	// 3. Get user
	u, err := s.userService.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if u.LevelID == nil || *u.LevelID != req.LevelID || u.StepID == nil || *u.StepID != req.StepID {
		return common.ErrBadRequest("You can only invest in your current level and step")
	}
	if u.CurrentStepCompleted {
		return common.ErrBadRequest("You have already completed the investment for this step. Please wait for the next step to be unlocked.")
	}

	// 5. Check balance
	if u.Balance < totalCost {
		return common.ErrBadRequest("Insufficient balance")
	}

	// 6. Deduct balance and add to withdrawable balance
	s.log.Infow("Processing investment balances", "user_id", userID, "total_cost", totalCost, "current_balance", u.Balance)
	if err := s.userService.UpdateBalance(ctx, userID, -totalCost); err != nil {
		return err
	}
	if err := s.userService.UpdateWithdrawableBalance(ctx, userID, totalCost); err != nil {
		// Rollback balance deduction
		s.userService.UpdateBalance(ctx, userID, totalCost)
		return err
	}

	// 7. Create transaction record
	tx := &Transaction{
		ID:     uuid.New().String(),
		UserID: userID,
		Type:   TypeInvestment,
		Amount: totalCost,
		Status: StatusCompleted,
		TxHash: common.GenerateHash(),
		Note:   fmt.Sprintf("Investment in Level %d Step %d (x%d sets)", req.LevelID, req.StepID, req.Quantity),
	}
	if err := s.repo.Create(ctx, tx); err != nil {
		// Rollback balances (manual)
		s.userService.UpdateBalance(ctx, userID, totalCost)
		s.userService.UpdateWithdrawableBalance(ctx, userID, -totalCost)
		return common.ErrInternal(err)
	}

	s.log.Infow("Investment successful, updating progress", "user_id", userID)

	// 8. Update user's progress to next step
	// Fetch all steps in this level
	steps, _, err := s.productService.ListStepsByLevel(ctx, req.LevelID, 100, 0)
	if err != nil {
		return nil // Investment successful, but progress tracking failed
	}

	var nextStepID uint
	foundCurrent := false
	for _, step := range steps {
		if foundCurrent {
			nextStepID = step.ID
			break
		}
		if step.ID == req.StepID {
			foundCurrent = true
		}
	}

	if nextStepID != 0 {
		// Move to next step
		completed := false
		reqUpdate := &user.UpdateUserRequest{
			StepID:               &nextStepID,
			CurrentStepCompleted: &completed,
		}
		s.log.Infow("Advancing user to next step", "user_id", userID, "next_step_id", nextStepID)
		if _, err := s.userService.Update(ctx, userID, reqUpdate); err != nil {
			s.log.Errorw("Failed to advance user to next step", "error", err, "user_id", userID)
		}
	} else {
		// Move to next level
		levels, _, err := s.productService.ListLevels(ctx, 100, 0)
		if err == nil {
			var nextLevelID uint
			foundLvl := false
			for _, lvl := range levels {
				if foundLvl {
					nextLevelID = lvl.ID
					break
				}
				if lvl.ID == req.LevelID {
					foundLvl = true
				}
			}

			if nextLevelID != 0 {
				// Get first step of next level
				firstSteps, _, err := s.productService.ListStepsByLevel(ctx, nextLevelID, 1, 0)
				var firstStepID *uint
				if err == nil && len(firstSteps) > 0 {
					id := firstSteps[0].ID
					firstStepID = &id
				}

				completed := false
				reqUpdate := &user.UpdateUserRequest{
					LevelID:              &nextLevelID,
					StepID:               firstStepID,
					CurrentStepCompleted: &completed,
				}
				s.log.Infow("Advancing user to next level", "user_id", userID, "next_level_id", nextLevelID, "first_step_id", firstStepID)
				if _, err := s.userService.Update(ctx, userID, reqUpdate); err != nil {
					s.log.Errorw("Failed to advance user to next level", "error", err, "user_id", userID)
				}
			} else {
				// NO MORE LEVELS - ABSOLUTE END
				s.log.Infow("User completed all levels and steps - marking current step as completed", "user_id", userID)
				completed := true
				reqUpdate := &user.UpdateUserRequest{
					CurrentStepCompleted: &completed,
				}
				if _, err := s.userService.Update(ctx, userID, reqUpdate); err != nil {
					s.log.Errorw("Failed to mark user step as completed", "error", err, "user_id", userID)
				}
			}
		}
	}

	return nil
}
