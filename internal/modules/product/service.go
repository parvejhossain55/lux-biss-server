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

func (s *ProductService) ListLevels(ctx context.Context) ([]*Level, error) {
	return s.repo.ListLevels(ctx)
}

func (s *ProductService) ListStepsByLevel(ctx context.Context, levelID uint) ([]*Step, error) {
	return s.repo.ListStepsByLevel(ctx, levelID)
}
