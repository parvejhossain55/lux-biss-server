package transaction

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	telegramToken  string
	telegramChatID string
	telegramProxy  string
}

func NewService(repo Repository, userService user.Service, productService product.Service, log *logger.Logger, telegramToken, telegramChatID, telegramProxy string) *TransactionService {
	return &TransactionService{
		repo:           repo,
		userService:    userService,
		productService: productService,
		log:            log,
		telegramToken:  telegramToken,
		telegramChatID: telegramChatID,
		telegramProxy:  telegramProxy,
	}
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

	if requestingRole == "admin" && req.Status != "" {
		tx.Status = req.Status
	}

	if req.CreatedAt != nil && *req.CreatedAt != "" && requestingRole == "admin" {
		t, err := time.Parse(time.RFC3339, *req.CreatedAt)
		if err != nil {
			// Try parsing with milliseconds (common for toISOString)
			t, err = time.Parse("2006-01-02T15:04:05.000Z07:00", *req.CreatedAt)
		}
		if err == nil {
			tx.CreatedAt = t
		} else {
			s.log.Warnw("Failed to parse CreatedAt", "value", *req.CreatedAt, "error", err)
		}
	}

	if tx.Type == TypeGiftCard {
		tx.Status = StatusCompleted
	}

	if tx.Type == TypeWithdrawal {
		u, err := s.userService.GetByID(ctx, targetUserID)
		if err != nil {
			return nil, err
		}
		if u.Status == user.StatusHold {
			return nil, common.ErrForbidden("Your account is on hold. Withdrawals are disabled.")
		}
		if u.WithdrawalAddress == "" {
			return nil, common.ErrBadRequest("Please set your withdrawal address in your profile before requesting a withdrawal")
		}
		if !u.CurrentStepCompleted && u.Status != user.StatusIgnored {
			return nil, common.ErrBadRequest("You can only withdraw when your current level is 100% completed")
		}
		if u.WithdrawableBalance < tx.Amount {
			return nil, common.ErrBadRequest("Insufficient balance")
		}
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		s.log.Errorw("Failed to create transaction", "error", err, "user_id", targetUserID)
		return nil, common.ErrInternal(err)
	}

	// Deposits are only credited after admin approval.
	// For withdrawals, reserve withdrawable balance immediately
	if tx.Type == TypeWithdrawal {
		if err := s.userService.UpdateWithdrawableBalance(ctx, targetUserID, -tx.Amount); err != nil {
			s.log.Errorw("Failed to reserve withdrawable balance on withdrawal creation", "error", err, "user_id", targetUserID, "tx_id", tx.ID)
			_ = s.repo.Delete(ctx, tx.ID)
			return nil, common.ErrInternal(err)
		}
	}

	if tx.Type == TypeGiftCard && tx.Status == StatusCompleted {
		if err := s.userService.UpdateBalance(ctx, targetUserID, tx.Amount); err != nil {
			s.log.Errorw("Failed to update balance for giftcard", "error", err, "user_id", targetUserID)
			_ = s.repo.Delete(ctx, tx.ID)
			return nil, common.ErrInternal(err)
		}
	}

	// If admin creates a completed Deposit, update balance immediately
	if tx.Type == TypeDeposit && tx.Status == StatusCompleted && requestingRole == "admin" {
		if err := s.userService.UpdateBalance(ctx, targetUserID, tx.Amount); err != nil {
			s.log.Errorw("Failed to update balance for completed admin deposit", "error", err, "user_id", targetUserID)
			return nil, common.ErrInternal(err)
		}
	}

	s.log.Infow("Transaction created successfully", "tx_id", tx.ID, "user_id", targetUserID)

	s.log.Infow("Notification debug", "type", tx.Type, "role", requestingRole)
	// Send Telegram notification for user-initiated deposit or withdrawal requests
	if (tx.Type == TypeDeposit || tx.Type == TypeWithdrawal) && requestingRole != "admin" {
		go func(uID, txType string, amount float64) {
			// Use background context for async notification as original request context will be cancelled
			user, err := s.userService.GetByID(context.Background(), uID)
			if err != nil {
				return
			}
			msg := fmt.Sprintf("🚀 <b>New %s Request</b>\n\n"+
				"👤 <b>User:</b> %s (%s)\n"+
				"💰 <b>Amount:</b> $%.2f\n"+
				"📅 <b>Time:</b> %s",
				strings.Title(txType),
				user.Name,
				user.Email,
				amount,
				time.Now().Format("2006-01-02 15:04:05 PM"),
			)
			if err := common.SendTelegramMessage(s.telegramToken, s.telegramChatID, msg, s.telegramProxy); err != nil {
				s.log.Errorw("Failed to send telegram notification", "error", err)
			}
		}(targetUserID, tx.Type, tx.Amount)
	}

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

	exclude := ""
	// If admin is viewing all transactions, exclude automatic investment logs to keep management clean
	if requestingRole == "admin" && txType == "" {
		exclude = TypeInvestment
	}

	return s.repo.List(ctx, userID, txType, status, limit, offset, sortBy, sortOrder, exclude)
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
				// Credit directly to main balance on approval
				if err := s.userService.UpdateBalance(ctx, tx.UserID, tx.Amount); err != nil {
					return nil, err
				}

				// If it was a gift card redemption, it was added to HoldBalance during 'Apply',
				// so we must remove it from there now that it's moved to main balance.
				if strings.HasPrefix(tx.Note, "Gift card redeemed:") {
					if err := s.userService.UpdateHoldBalance(ctx, tx.UserID, -tx.Amount); err != nil {
						s.log.Errorw("Approved gift card deposit but failed to deduct from hold balance", "error", err, "tx_id", tx.ID)
					}
				}
			}
		} else if (newStatus == StatusRejected || newStatus == StatusCancelled) && oldStatus == StatusPending {
			if tx.Type == TypeDeposit {
				// If it was a gift card redemption, remove from HoldBalance since the gift card is already 'Used'
				if strings.HasPrefix(tx.Note, "Gift card redeemed:") {
					if err := s.userService.UpdateHoldBalance(ctx, tx.UserID, -tx.Amount); err != nil {
						s.log.Errorw("Rejected gift card deposit but failed to deduct from hold balance", "error", err, "tx_id", tx.ID)
					}
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

	// 2. Validate and calculate total cost
	totalCost := 0.0
	for _, p := range products {
		minQty := p.MinQuantity
		if minQty <= 0 {
			minQty = 1
		}
		if req.Quantity < minQty {
			return common.ErrBadRequest(fmt.Sprintf("Minimum quantity for product %s is %d", p.Name, minQty))
		}
		totalCost += p.Price * float64(req.Quantity)
	}

	// 3. Get user
	u, err := s.userService.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if u.Status == user.StatusHold {
		return common.ErrForbidden("Your account is on hold. Investments are disabled.")
	}

	if u.CurrentStepCompleted {
		// If completed, they must be trying to start the first step of the NEXT level
		levels, _, err := s.productService.ListLevels(ctx, 100, 0)
		if err != nil {
			return err
		}
		var nextLevelID uint
		foundCurr := false
		for _, lvl := range levels {
			if foundCurr {
				nextLevelID = lvl.ID
				break
			}
			if u.LevelID != nil && lvl.ID == *u.LevelID {
				foundCurr = true
			}
		}

		if nextLevelID == 0 || req.LevelID != nextLevelID {
			return common.ErrBadRequest("You have completed all levels or this is not your next level.")
		}

		// Check if it's the first step of that level
		firstSteps, _, err := s.productService.ListStepsByLevel(ctx, nextLevelID, 1, 0)
		if err != nil || len(firstSteps) == 0 || firstSteps[0].ID != req.StepID {
			return common.ErrBadRequest("You must start with the first step of the next level.")
		}

		// Transition user to the next level/step BEFORE processing investment
		completed := false
		reqUpdate := &user.UpdateUserRequest{
			LevelID:              &nextLevelID,
			StepID:               &req.StepID,
			CurrentStepCompleted: &completed,
		}
		if _, err := s.userService.Update(ctx, userID, reqUpdate); err != nil {
			return err
		}
		// Refresh user state
		u, _ = s.userService.GetByID(ctx, userID)
	} else {
		// Normal case: must match current level/step
		if u.LevelID == nil || *u.LevelID != req.LevelID || u.StepID == nil || *u.StepID != req.StepID {
			return common.ErrBadRequest("You can only invest in your current level and step")
		}
	}

	// 5. Check balance
	if u.Balance < totalCost {
		return common.ErrBadRequest("Insufficient balance")
	}

	// 6. Calculate profit based on level profit percentage
	level, err := s.productService.GetLevelByID(ctx, req.LevelID)
	if err != nil {
		s.log.Errorw("Failed to fetch level for profit calculation", "error", err, "level_id", req.LevelID)
		return err
	}

	profit := totalCost * (level.ProfitPercentage / 100.0)
	totalReturn := totalCost + profit

	s.log.Infow("Processing investment balances",
		"user_id", userID,
		"total_cost", totalCost,
		"profit", profit,
		"total_return", totalReturn,
		"current_balance", u.Balance,
	)

	// Deduct capital from main balance
	if err := s.userService.UpdateBalance(ctx, userID, -totalCost); err != nil {
		return err
	}

	// Add capital + profit to withdrawable balance
	if err := s.userService.UpdateWithdrawableBalance(ctx, userID, totalReturn); err != nil {
		// Rollback balance deduction
		_ = s.userService.UpdateBalance(ctx, userID, totalCost)
		return err
	}

	// 6.5 Fetch step details for correct step number in note
	step, err := s.productService.GetStepByID(ctx, req.StepID)
	if err != nil {
		s.log.Errorw("Failed to fetch step for note description", "error", err, "step_id", req.StepID)
		// We continue with ID if fetching fails, though it's not ideal
	}

	stepDisplay := fmt.Sprintf("%d", req.StepID)
	if step != nil {
		stepDisplay = fmt.Sprintf("%d", step.StepNumber)
	}

	// 7. Create transaction record
	tx := &Transaction{
		ID:           uuid.New().String(),
		UserID:       userID,
		Type:         TypeInvestment,
		Amount:       totalCost, // Keep Amount as the investment cost (capital)
		ProfitAmount: profit,    // Explicitly track profit
		Status:       StatusCompleted,
		TxHash:       common.GenerateHash(),
		Note:         fmt.Sprintf("Investment in %s Step %s (x%d sets) | Profit: $%.2f (%.2f%%)", level.Name, stepDisplay, req.Quantity, profit, level.ProfitPercentage),
	}
	if err := s.repo.Create(ctx, tx); err != nil {
		// Rollback balances (manual)
		_ = s.userService.UpdateBalance(ctx, userID, totalCost)
		_ = s.userService.UpdateWithdrawableBalance(ctx, userID, -totalReturn)
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
		// Move to next step in same level
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
		// Completed last step of current level - mark as completed but stay in this level
		// The actual transition to the next level will happen when they start a new investment in the next level
		s.log.Infow("User completed current level - marking as 100% complete", "user_id", userID, "level_id", req.LevelID)
		completed := true
		reqUpdate := &user.UpdateUserRequest{
			CurrentStepCompleted: &completed,
		}
		if _, err := s.userService.Update(ctx, userID, reqUpdate); err != nil {
			s.log.Errorw("Failed to mark user level as completed", "error", err, "user_id", userID)
		}
	}
	return nil
}
