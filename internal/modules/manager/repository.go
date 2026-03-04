package manager

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

func (r *GormRepository) Create(ctx context.Context, manager *Manager) error {
	manager.ID = uuid.New().String()

	result := r.db.WithContext(ctx).Create(manager)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id string) (*Manager, error) {
	var manager Manager
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&manager)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("Manager")
		}
		return nil, result.Error
	}

	return &manager, nil
}

func (r *GormRepository) List(ctx context.Context, limit, offset int) ([]*Manager, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&Manager{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var managers []*Manager
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&managers)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return managers, total, nil
}

func (r *GormRepository) Update(ctx context.Context, manager *Manager) error {
	result := r.db.WithContext(ctx).
		Model(manager).
		Updates(map[string]interface{}{
			"name":              manager.Name,
			"telegram_username": manager.TelegramUsername,
			"telegram_link":     manager.TelegramLink,
			"profile_photo":     manager.ProfilePhoto,
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrNotFound("Manager")
	}

	return nil
}

func (r *GormRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&Manager{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrNotFound("Manager")
	}

	return nil
}
