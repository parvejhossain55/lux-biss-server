package manager

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Manager struct {
	ID               string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name             string         `json:"name" gorm:"type:varchar(100);not null"`
	TelegramUsername string         `json:"telegram_username" gorm:"type:varchar(100)"`
	TelegramLink     string         `json:"telegram_link" gorm:"type:varchar(255)"`
	ProfilePhoto     string         `json:"profile_photo" gorm:"type:varchar(255)"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Manager) TableName() string {
	return "managers"
}

type Repository interface {
	Create(ctx context.Context, manager *Manager) error
	GetByID(ctx context.Context, id string) (*Manager, error)
	List(ctx context.Context, limit, offset int) ([]*Manager, int64, error)
	Update(ctx context.Context, manager *Manager) error
	Delete(ctx context.Context, id string) error
}

type Service interface {
	Create(ctx context.Context, req *CreateManagerRequest) (*Manager, error)
	GetByID(ctx context.Context, id string) (*Manager, error)
	List(ctx context.Context, limit, offset int) ([]*Manager, int64, error)
	Update(ctx context.Context, id string, req *UpdateManagerRequest) (*Manager, error)
	Delete(ctx context.Context, id string) error
}
