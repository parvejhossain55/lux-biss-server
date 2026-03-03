package user

import (
	"context"
	"strings"

	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/logger"
	"github.com/parvej/luxbiss_server/pkg/hash"
)

type UserService struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) *UserService {
	return &UserService{repo: repo, log: log}
}

func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
	existing, _ := s.repo.GetByEmail(ctx, strings.ToLower(req.Email))
	if existing != nil {
		return nil, common.ErrConflict("A user with this email already exists")
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		s.log.Errorw("Failed to hash password", "error", err)
		return nil, common.ErrInternal(err)
	}

	user := &User{
		Name:             req.Name,
		Email:            strings.ToLower(req.Email),
		Password:         hashedPassword,
		Role:             req.Role,
		ProfilePhoto:     req.ProfilePhoto,
		TelegramUsername: req.TelegramUsername,
		TelegramLink:     req.TelegramLink,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.log.Errorw("Failed to create user", "error", err, "email", req.Email)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("User created successfully", "user_id", user.ID, "email", user.Email)
	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.repo.GetByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) List(ctx context.Context, limit, offset int) ([]*User, int64, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *UserService) Update(ctx context.Context, id string, req *UpdateUserRequest) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = strings.ToLower(req.Email)
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.ProfilePhoto != nil {
		user.ProfilePhoto = *req.ProfilePhoto
	}
	if req.TelegramUsername != nil {
		user.TelegramUsername = *req.TelegramUsername
	}
	if req.TelegramLink != nil {
		user.TelegramLink = *req.TelegramLink
	}

	if err := s.repo.Update(ctx, user); err != nil {
		s.log.Errorw("Failed to update user", "error", err, "user_id", id)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("User updated successfully", "user_id", id)
	return user, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	if err := s.repo.UpdatePassword(ctx, id, hashedPassword); err != nil {
		s.log.Errorw("Failed to update user password", "error", err, "user_id", id)
		return common.ErrInternal(err)
	}
	s.log.Infow("User password updated successfully", "user_id", id)
	return nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.log.Infow("User deleted successfully", "user_id", id)
	return nil
}
