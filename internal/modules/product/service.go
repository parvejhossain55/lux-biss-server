package product

import (
	"context"

	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/logger"
)

type ProductService struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) *ProductService {
	return &ProductService{repo: repo, log: log}
}

func (s *ProductService) Create(ctx context.Context, req *CreateProductRequest) (*Product, error) {
	product := &Product{
		LevelID:     req.LevelID,
		StepID:      req.StepID,
		Name:        req.Name,
		Price:       req.Price,
		Rating:      req.Rating,
		MinQuantity: req.MinQuantity,
		MaxQuantity: req.MaxQuantity,
		ImageURL:    req.ImageURL,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		s.log.Errorw("Failed to create product", "error", err, "name", req.Name)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("Product created successfully", "product_id", product.ID, "name", product.Name)
	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) List(ctx context.Context, limit, offset int, sortBy, order string, levelID, stepID uint) ([]*Product, int64, error) {
	return s.repo.List(ctx, limit, offset, sortBy, order, levelID, stepID)
}

func (s *ProductService) Update(ctx context.Context, id string, req *UpdateProductRequest) (*Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.LevelID != nil {
		product.LevelID = *req.LevelID
	}
	if req.StepID != nil {
		product.StepID = *req.StepID
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Rating != nil {
		product.Rating = *req.Rating
	}
	if req.MinQuantity != nil {
		product.MinQuantity = *req.MinQuantity
	}
	if req.MaxQuantity != nil {
		product.MaxQuantity = *req.MaxQuantity
	}
	if req.ImageURL != nil {
		product.ImageURL = *req.ImageURL
	}
	if req.Description != nil {
		product.Description = *req.Description
	}

	if err := s.repo.Update(ctx, product); err != nil {
		s.log.Errorw("Failed to update product", "error", err, "product_id", id)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("Product updated successfully", "product_id", id)
	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.log.Infow("Product deleted successfully", "product_id", id)
	return nil
}

// Level operations
func (s *ProductService) CreateLevel(ctx context.Context, req *CreateLevelRequest) (*Level, error) {
	level := &Level{
		Name: req.Name,
	}
	if err := s.repo.CreateLevel(ctx, level); err != nil {
		s.log.Errorw("Failed to create level", "error", err, "name", req.Name)
		return nil, common.ErrInternal(err)
	}
	s.log.Infow("Level created successfully", "level_id", level.ID, "name", level.Name)
	return level, nil
}

func (s *ProductService) GetLevelByID(ctx context.Context, id uint) (*Level, error) {
	return s.repo.GetLevelByID(ctx, id)
}

func (s *ProductService) ListLevels(ctx context.Context, limit, offset int) ([]*Level, int64, error) {
	return s.repo.ListLevels(ctx, limit, offset)
}

func (s *ProductService) UpdateLevel(ctx context.Context, id uint, req *UpdateLevelRequest) (*Level, error) {
	level, err := s.repo.GetLevelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		level.Name = *req.Name
	}
	if err := s.repo.UpdateLevel(ctx, level); err != nil {
		s.log.Errorw("Failed to update level", "error", err, "level_id", id)
		return nil, common.ErrInternal(err)
	}
	return level, nil
}

func (s *ProductService) DeleteLevel(ctx context.Context, id uint) error {
	return s.repo.DeleteLevel(ctx, id)
}

// Step operations
func (s *ProductService) CreateStep(ctx context.Context, req *CreateStepRequest) (*Step, error) {
	step := &Step{
		LevelID:    req.LevelID,
		StepNumber: req.StepNumber,
		Name:       req.Name,
	}
	if err := s.repo.CreateStep(ctx, step); err != nil {
		s.log.Errorw("Failed to create step", "error", err, "level_id", req.LevelID)
		return nil, common.ErrInternal(err)
	}
	s.log.Infow("Step created successfully", "step_id", step.ID, "level_id", step.LevelID)
	return step, nil
}

func (s *ProductService) GetStepByID(ctx context.Context, id uint) (*Step, error) {
	return s.repo.GetStepByID(ctx, id)
}

func (s *ProductService) ListStepsByLevel(ctx context.Context, levelID uint, limit, offset int) ([]*Step, int64, error) {
	return s.repo.ListStepsByLevel(ctx, levelID, limit, offset)
}

func (s *ProductService) UpdateStep(ctx context.Context, id uint, req *UpdateStepRequest) (*Step, error) {
	step, err := s.repo.GetStepByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.LevelID != nil {
		step.LevelID = *req.LevelID
	}
	if req.StepNumber != nil {
		step.StepNumber = *req.StepNumber
	}
	if req.Name != nil {
		step.Name = *req.Name
	}
	if err := s.repo.UpdateStep(ctx, step); err != nil {
		s.log.Errorw("Failed to update step", "error", err, "step_id", id)
		return nil, common.ErrInternal(err)
	}
	return step, nil
}

func (s *ProductService) DeleteStep(ctx context.Context, id uint) error {
	return s.repo.DeleteStep(ctx, id)
}
