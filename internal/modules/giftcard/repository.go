package giftcard

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

func (r *GormRepository) Create(ctx context.Context, giftcard *Giftcard) error {
	giftcard.ID = uuid.New().String()
	if giftcard.Status == "" {
		giftcard.Status = StatusAvailable
	}

	result := r.db.WithContext(ctx).Create(giftcard)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id string) (*Giftcard, error) {
	var giftcard Giftcard
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&giftcard)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("Giftcard")
		}
		return nil, result.Error
	}
	return &giftcard, nil
}

func (r *GormRepository) GetByCode(ctx context.Context, code string) (*Giftcard, error) {
	var giftcard Giftcard
	result := r.db.WithContext(ctx).Where("redeem_code = ?", code).First(&giftcard)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("Giftcard")
		}
		return nil, result.Error
	}
	return &giftcard, nil
}

func (r *GormRepository) List(ctx context.Context, limit, offset int) ([]*Giftcard, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&Giftcard{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var giftcards []*Giftcard
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&giftcards)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return giftcards, total, nil
}

func (r *GormRepository) Update(ctx context.Context, giftcard *Giftcard) error {
	result := r.db.WithContext(ctx).Save(giftcard)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound("Giftcard")
	}
	return nil
}

func (r *GormRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&Giftcard{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound("Giftcard")
	}
	return nil
}
