package manager

import (
	"context"

	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/logger"
)

type ManagerService struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) *ManagerService {
	return &ManagerService{repo: repo, log: log}
}

func (s *ManagerService) Create(ctx context.Context, req *CreateManagerRequest) (*Manager, error) {
	manager := &Manager{
		Name:             req.Name,
		TelegramUsername: req.TelegramUsername,
		ProfilePhoto:     req.ProfilePhoto,
	}

	if err := s.repo.Create(ctx, manager); err != nil {
		s.log.Errorw("Failed to create manager", "error", err, "name", req.Name)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("Manager created successfully", "manager_id", manager.ID)
	return manager, nil
}

func (s *ManagerService) GetByID(ctx context.Context, id string) (*Manager, error) {
	manager, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return manager, nil
}

func (s *ManagerService) List(ctx context.Context, limit, offset int) ([]*Manager, int64, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *ManagerService) Update(ctx context.Context, id string, req *UpdateManagerRequest) (*Manager, error) {
	manager, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		manager.Name = *req.Name
	}
	if req.TelegramUsername != nil {
		manager.TelegramUsername = *req.TelegramUsername
	}

	if req.ProfilePhoto != nil {
		manager.ProfilePhoto = *req.ProfilePhoto
	}

	if err := s.repo.Update(ctx, manager); err != nil {
		s.log.Errorw("Failed to update manager", "error", err, "manager_id", id)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("Manager updated successfully", "manager_id", id)
	return manager, nil
}

func (s *ManagerService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.log.Infow("Manager deleted successfully", "manager_id", id)
	return nil
}
