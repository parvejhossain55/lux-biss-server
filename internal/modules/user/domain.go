package user

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID               string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name             string         `json:"name" gorm:"type:varchar(100);not null"`
	Email            string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password         string         `json:"-" gorm:"type:varchar(255);not null"`
	Role             string         `json:"role" gorm:"type:varchar(20);not null;default:'user';index"`
	IsActive         bool           `json:"is_active" gorm:"not null;default:true"`
	ProfilePhoto     string         `json:"profile_photo" gorm:"type:varchar(255)"`
	TelegramUsername string         `json:"telegram_username" gorm:"type:varchar(255)"`
	TelegramLink     string         `json:"telegram_link" gorm:"type:varchar(255)"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, limit, offset int) ([]*User, int64, error)
	Update(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, id string, hashedPassword string) error
	Delete(ctx context.Context, id string) error
}

type Service interface {
	Create(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, limit, offset int) ([]*User, int64, error)
	Update(ctx context.Context, id string, req *UpdateUserRequest) (*User, error)
	UpdatePassword(ctx context.Context, id string, hashedPassword string) error
	Delete(ctx context.Context, id string) error
}
