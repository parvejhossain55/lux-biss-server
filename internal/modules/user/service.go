package user

import (
	"context"
	"strings"

	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/logger"
	"github.com/parvej/luxbiss_server/pkg/hash"
)

type UserService struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) *UserService {
	return &UserService{repo: repo, log: log}
}

func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
	existing, _ := s.repo.GetByEmail(ctx, strings.ToLower(req.Email))
	if existing != nil {
		return nil, common.ErrConflict("A user with this email already exists")
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		s.log.Errorw("Failed to hash password", "error", err)
		return nil, common.ErrInternal(err)
	}

	user := &User{
		Name:         req.Name,
		Email:        strings.ToLower(req.Email),
		Password:     hashedPassword,
		Role:         req.Role,
		ProfilePhoto: req.ProfilePhoto,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.log.Errorw("Failed to create user", "error", err, "email", req.Email)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("User created successfully", "user_id", user.ID, "email", user.Email)
	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.repo.GetByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) List(ctx context.Context, status string, limit, offset int) ([]*User, int64, error) {
	return s.repo.List(ctx, status, limit, offset)
}

func (s *UserService) Update(ctx context.Context, id string, req *UpdateUserRequest) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = strings.ToLower(req.Email)
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.ProfilePhoto != nil {
		user.ProfilePhoto = *req.ProfilePhoto
	}
	if req.Balance != nil {
		user.Balance = *req.Balance
	}
	if req.HoldBalance != nil {
		user.HoldBalance = *req.HoldBalance
	}
	if req.WithdrawableBalance != nil {
		user.WithdrawableBalance = *req.WithdrawableBalance
	}
	if req.DateOfBirth != nil {
		user.DateOfBirth = *req.DateOfBirth
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Address != nil {
		user.Address = *req.Address
	}
	if req.Country != nil {
		user.Country = *req.Country
	}
	if req.PaymentMethod != nil {
		user.PaymentMethod = *req.PaymentMethod
	}
	if req.PaymentCurrency != nil {
		user.PaymentCurrency = *req.PaymentCurrency
	}
	if req.PaymentNetwork != nil {
		user.PaymentNetwork = *req.PaymentNetwork
	}
	if req.WithdrawalAddress != nil {
		user.WithdrawalAddress = *req.WithdrawalAddress
	}
	if req.LevelID != nil {
		user.LevelID = req.LevelID
	}
	if req.StepID != nil {
		user.StepID = req.StepID
	}
	if req.CurrentStepCompleted != nil {
		user.CurrentStepCompleted = *req.CurrentStepCompleted
	}

	if err := s.repo.Update(ctx, user); err != nil {
		s.log.Errorw("Failed to update user", "error", err, "user_id", id)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("User updated successfully", "user_id", id)
	return user, nil
}

func (s *UserService) UpdateBalance(ctx context.Context, userID string, amount float64) error {
	if err := s.repo.UpdateBalance(ctx, userID, amount); err != nil {
		s.log.Errorw("Failed to update user balance", "error", err, "user_id", userID, "amount", amount)
		return common.ErrInternal(err)
	}
	s.log.Infow("User balance updated successfully", "user_id", userID, "amount", amount)
	return nil
}

func (s *UserService) UpdateHoldBalance(ctx context.Context, userID string, amount float64) error {
	if err := s.repo.UpdateHoldBalance(ctx, userID, amount); err != nil {
		s.log.Errorw("Failed to update user hold balance", "error", err, "user_id", userID, "amount", amount)
		return common.ErrInternal(err)
	}
	s.log.Infow("User hold balance updated successfully", "user_id", userID, "amount", amount)
	return nil
}

func (s *UserService) UpdateWithdrawableBalance(ctx context.Context, userID string, amount float64) error {
	if err := s.repo.UpdateWithdrawableBalance(ctx, userID, amount); err != nil {
		s.log.Errorw("Failed to update user withdrawable balance", "error", err, "user_id", userID, "amount", amount)
		return common.ErrInternal(err)
	}
	s.log.Infow("User withdrawable balance updated successfully", "user_id", userID, "amount", amount)
	return nil
}

func (s *UserService) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	if err := s.repo.UpdatePassword(ctx, id, hashedPassword); err != nil {
		s.log.Errorw("Failed to update user password", "error", err, "user_id", id)
		return common.ErrInternal(err)
	}
	s.log.Infow("User password updated successfully", "user_id", id)
	return nil
}

func (s *UserService) ApproveHoldBalance(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.HoldBalance <= 0 {
		return nil, common.ErrBadRequest("User has no hold balance to approve")
	}

	amount := user.HoldBalance
	user.Balance += amount
	user.HoldBalance = 0

	if err := s.repo.Update(ctx, user); err != nil {
		s.log.Errorw("Failed to approve hold balance", "error", err, "user_id", id)
		return nil, common.ErrInternal(err)
	}

	// Sync transactions
	if err := s.CompletePendingTransactions(ctx, id); err != nil {
		s.log.Errorw("Failed to sync transactions on hold balance approval", "error", err, "user_id", id)
		// We don't return error here because the balance itself was already updated
	}

	s.log.Infow("User hold balance approved successfully", "user_id", id, "amount", amount)
	return user, nil
}

func (s *UserService) CompletePendingTransactions(ctx context.Context, userID string) error {
	return s.repo.CompletePendingTransactions(ctx, userID)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.log.Infow("User deleted successfully", "user_id", id)
	return nil
}

func (s *UserService) AdvanceUsersToNextStep(ctx context.Context, levelID, currentStepID, nextLevelID, nextStepID uint) error {
	if err := s.repo.AdvanceUsersToNextStep(ctx, levelID, currentStepID, nextLevelID, nextStepID); err != nil {
		s.log.Errorw("Failed to advance users to next step", "error", err, "level_id", levelID, "current_step_id", currentStepID, "next_step_id", nextStepID, "next_level_id", nextLevelID)
		return common.ErrInternal(err)
	}
	s.log.Infow("Advanced users to next step", "level_id", levelID, "current_step_id", currentStepID, "next_step_id", nextStepID, "next_level_id", nextLevelID)
	return nil
}
