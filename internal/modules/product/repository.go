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

func (r *GormRepository) List(ctx context.Context, limit, offset int) ([]*Product, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []*Product
	result := r.db.WithContext(ctx).
		Preload("Level").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products)
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
