package user

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

func (r *GormRepository) Create(ctx context.Context, user *User) error {
	user.ID = uuid.New().String()
	user.IsActive = true

	if user.Role == "" {
		user.Role = RoleUser
	}

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
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
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
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
		Updates(map[string]interface{}{
			"name":               user.Name,
			"email":              user.Email,
			"role":               user.Role,
			"is_active":          user.IsActive,
			"profile_photo":      user.ProfilePhoto,
			"telegram_username":  user.TelegramUsername,
			"telegram_link":      user.TelegramLink,
			"date_of_birth":      user.DateOfBirth,
			"gender":             user.Gender,
			"phone":              user.Phone,
			"address":            user.Address,
			"country":            user.Country,
			"payment_method":     user.PaymentMethod,
			"payment_currency":   user.PaymentCurrency,
			"payment_network":    user.PaymentNetwork,
			"withdrawal_address": user.WithdrawalAddress,
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
