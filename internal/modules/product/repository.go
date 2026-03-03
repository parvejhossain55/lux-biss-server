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

func (r *GormRepository) ListLevels(ctx context.Context) ([]*Level, error) {
	var levels []*Level
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&levels).Error; err != nil {
		return nil, err
	}
	return levels, nil
}

func (r *GormRepository) ListStepsByLevel(ctx context.Context, levelID uint) ([]*Step, error) {
	var steps []*Step
	if err := r.db.WithContext(ctx).Where("level_id = ?", levelID).Order("step_number ASC").Find(&steps).Error; err != nil {
		return nil, err
	}
	return steps, nil
}
