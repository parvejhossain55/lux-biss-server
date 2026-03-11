package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/parvej/luxbiss_server/internal/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, user *User) error {
	// If a user with the same email was soft-deleted, hard-delete it to permit recreation/re-registration
	if err := r.db.WithContext(ctx).Unscoped().Where("email = ? AND deleted_at IS NOT NULL", user.Email).Delete(&User{}).Error; err != nil {
		return err
	}

	user.ID = uuid.New().String()
	user.Status = StatusActive

	// if user.Role == "" {
	user.Role = RoleUser
	// }

	// Assign a random manager to the new user
	var managerID string
	// 'RANDOM()' for PostgreSQL / SQLite.
	err := r.db.WithContext(ctx).Table("managers").Select("id").Order("RANDOM()").Limit(1).Scan(&managerID).Error
	if err == nil && managerID != "" {
		user.ManagerID = &managerID
	}

	// Assign default level (lowest ID) and step (lowest step_number for that level)
	var levelID uint
	err = r.db.WithContext(ctx).Table("levels").Select("id").Order("id ASC").Limit(1).Scan(&levelID).Error
	if err == nil && levelID != 0 {
		user.LevelID = &levelID
		var stepID uint
		err = r.db.WithContext(ctx).Table("steps").Select("id").Where("level_id = ?", levelID).Order("step_number ASC").Limit(1).Scan(&stepID).Error
		if err == nil && stepID != 0 {
			user.StepID = &stepID
		}
	}

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).Preload("Manager").Preload("Level").Preload("Step").Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("User")
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *GormRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).Preload("Manager").Preload("Level").Preload("Step").Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("User")
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *GormRepository) List(ctx context.Context, limit, offset int) ([]*User, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []*User
	result := r.db.WithContext(ctx).
		Preload("Manager").
		Preload("Level").
		Preload("Step").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

func (r *GormRepository) Update(ctx context.Context, user *User) error {
	result := r.db.WithContext(ctx).
		Model(user).
		Omit(clause.Associations).
		Updates(map[string]interface{}{
			"name":                   user.Name,
			"email":                  user.Email,
			"role":                   user.Role,
			"status":                 user.Status,
			"profile_photo":          user.ProfilePhoto,
			"date_of_birth":          user.DateOfBirth,
			"gender":                 user.Gender,
			"phone":                  user.Phone,
			"address":                user.Address,
			"country":                user.Country,
			"payment_method":         user.PaymentMethod,
			"payment_currency":       user.PaymentCurrency,
			"payment_network":        user.PaymentNetwork,
			"withdrawal_address":     user.WithdrawalAddress,
			"balance":                user.Balance,
			"hold_balance":           user.HoldBalance,
			"withdrawable_balance":   user.WithdrawableBalance,
			"level_id":               user.LevelID,
			"step_id":                user.StepID,
			"current_step_completed": user.CurrentStepCompleted,
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrNotFound("User")
	}

	return nil
}

func (r *GormRepository) UpdateBalance(ctx context.Context, userID string, amount float64) error {
	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound("User")
	}
	return nil
}

func (r *GormRepository) UpdateHoldBalance(ctx context.Context, userID string, amount float64) error {
	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Update("hold_balance", gorm.Expr("hold_balance + ?", amount))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound("User")
	}
	return nil
}

func (r *GormRepository) UpdateWithdrawableBalance(ctx context.Context, userID string, amount float64) error {
	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Update("withdrawable_balance", gorm.Expr("withdrawable_balance + ?", amount))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound("User")
	}
	return nil
}

func (r *GormRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Update("password", hashedPassword)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrNotFound("User")
	}

	return nil
}

func (r *GormRepository) CompletePendingTransactions(ctx context.Context, userID string) error {
	// Execute raw SQL as transaction module is separate but we want to mark all its pending deposits as completed
	// We only complete "deposit" types that are "processing"
	result := r.db.WithContext(ctx).Exec(
		"UPDATE transactions SET status = 'completed', updated_at = NOW() WHERE user_id = ? AND status = 'processing' AND type = 'deposit'",
		userID,
	)
	return result.Error
}

func (r *GormRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrNotFound("User")
	}

	return nil
}

func (r *GormRepository) AdvanceUsersToNextStep(ctx context.Context, levelID, currentStepID, nextLevelID, nextStepID uint) error {
	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("level_id = ? AND step_id = ?", levelID, currentStepID).
		Updates(map[string]interface{}{
			"level_id":               nextLevelID,
			"step_id":                nextStepID,
			"current_step_completed": false,
		}).Error

	return result
}
