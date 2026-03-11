package product

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

func (r *GormRepository) Create(ctx context.Context, product *Product) error {
	product.ID = uuid.New().String()

	result := r.db.WithContext(ctx).Create(product)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	var product Product
	result := r.db.WithContext(ctx).Preload("Level").Where("id = ?", id).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("Product")
		}
		return nil, result.Error
	}

	return &product, nil
}

func (r *GormRepository) List(ctx context.Context, limit, offset int, sortBy, order string, levelID, stepID uint) ([]*Product, int64, error) {
	query := r.db.WithContext(ctx).Model(&Product{})

	// Apply filtering
	if levelID > 0 {
		query = query.Where("level_id = ?", levelID)
	}
	if stepID > 0 {
		query = query.Where("step_id = ?", stepID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Preload("Level").Preload("Step")

	// Apply sorting
	if sortBy != "" {
		if order == "" {
			order = "asc"
		}
		query = query.Order(sortBy + " " + order)
	} else {
		query = query.Order("created_at DESC")
	}

	var products []*Product
	result := query.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return products, total, nil
}

func (r *GormRepository) Update(ctx context.Context, product *Product) error {
	result := r.db.WithContext(ctx).
		Model(product).
		Where("id = ?", product.ID).
		Updates(product)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrNotFound("Product")
	}

	return nil
}

func (r *GormRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&Product{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrNotFound("Product")
	}

	return nil
}

// Level operations
func (r *GormRepository) CreateLevel(ctx context.Context, level *Level) error {
	// If a level with the same name was soft-deleted, clean up EVERYTHING associated with it to permit recreation
	var oldLevelIDs []uint
	if err := r.db.WithContext(ctx).Unscoped().Model(&Level{}).Where("name = ? AND deleted_at IS NOT NULL", level.Name).Pluck("id", &oldLevelIDs).Error; err == nil && len(oldLevelIDs) > 0 {
		_ = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Hard delete all associated records first
			tx.Unscoped().Where("level_id IN ?", oldLevelIDs).Delete(&Product{})
			tx.Unscoped().Where("level_id IN ?", oldLevelIDs).Delete(&Step{})
			tx.Unscoped().Where("id IN ?", oldLevelIDs).Delete(&Level{})
			return nil
		})
	}

	return r.db.WithContext(ctx).Create(level).Error
}

func (r *GormRepository) GetLevelByID(ctx context.Context, id uint) (*Level, error) {
	var level Level
	if err := r.db.WithContext(ctx).First(&level, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("Level")
		}
		return nil, err
	}
	return &level, nil
}

func (r *GormRepository) ListLevels(ctx context.Context, limit, offset int) ([]*Level, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&Level{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var levels []*Level
	query := r.db.WithContext(ctx).Model(&Level{}).Order("id ASC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&levels).Error; err != nil {
		return nil, 0, err
	}
	return levels, total, nil
}

func (r *GormRepository) UpdateLevel(ctx context.Context, level *Level) error {
	result := r.db.WithContext(ctx).Save(level)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound("Level")
	}
	return nil
}

func (r *GormRepository) DeleteLevel(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete products associated with this level
		if err := tx.Delete(&Product{}, "level_id = ?", id).Error; err != nil {
			return err
		}

		// Delete steps associated with this level
		if err := tx.Delete(&Step{}, "level_id = ?", id).Error; err != nil {
			return err
		}

		// Delete the level itself
		result := tx.Delete(&Level{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return common.ErrNotFound("Level")
		}
		return nil
	})
}

// Step operations
func (r *GormRepository) CreateStep(ctx context.Context, step *Step) error {
	return r.db.WithContext(ctx).Create(step).Error
}

func (r *GormRepository) GetStepByID(ctx context.Context, id uint) (*Step, error) {
	var step Step
	if err := r.db.WithContext(ctx).Preload("Level").First(&step, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("Step")
		}
		return nil, err
	}
	return &step, nil
}

func (r *GormRepository) ListStepsByLevel(ctx context.Context, levelID uint, limit, offset int) ([]*Step, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&Step{}).Where("level_id = ?", levelID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var steps []*Step
	query := r.db.WithContext(ctx).Model(&Step{}).Where("level_id = ?", levelID).Order("step_number ASC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&steps).Error; err != nil {
		return nil, 0, err
	}
	return steps, total, nil
}

func (r *GormRepository) UpdateStep(ctx context.Context, step *Step) error {
	result := r.db.WithContext(ctx).Save(step)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound("Step")
	}
	return nil
}

func (r *GormRepository) DeleteStep(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete products associated with this step
		if err := tx.Delete(&Product{}, "step_id = ?", id).Error; err != nil {
			return err
		}

		// Delete the step itself
		result := tx.Delete(&Step{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return common.ErrNotFound("Step")
		}
		return nil
	})
}
