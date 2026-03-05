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
	ID               string  `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name             string  `json:"name" gorm:"type:varchar(100);not null"`
	Email            string  `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password         string  `json:"-" gorm:"type:varchar(255);not null"`
	Role             string  `json:"role" gorm:"type:varchar(20);not null;default:'user';index"`
	IsActive         bool    `json:"is_active" gorm:"not null;default:true"`
	ProfilePhoto     string  `json:"profile_photo" gorm:"type:varchar(255)"`
	TelegramUsername string  `json:"telegram_username" gorm:"type:varchar(255)"`
	TelegramLink     string  `json:"telegram_link" gorm:"type:varchar(255)"`
	Balance          float64 `json:"balance" gorm:"type:decimal(15,2);not null;default:0"`
	// Personal Information
	DateOfBirth string `json:"date_of_birth" gorm:"type:varchar(20)"`
	Gender      string `json:"gender" gorm:"type:varchar(20)"`
	Phone       string `json:"phone" gorm:"type:varchar(30)"`
	Address     string `json:"address" gorm:"type:varchar(500)"`
	Country     string `json:"country" gorm:"type:varchar(100)"`
	// Payment Wallet Information
	PaymentMethod     string         `json:"payment_method" gorm:"type:varchar(100)"`
	PaymentCurrency   string         `json:"payment_currency" gorm:"type:varchar(50)"`
	PaymentNetwork    string         `json:"payment_network" gorm:"type:varchar(100)"`
	WithdrawalAddress string         `json:"withdrawal_address" gorm:"type:varchar(255)"`
	CreatedAt         time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
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
	UpdateBalance(ctx context.Context, userID string, amount float64) error
	UpdatePassword(ctx context.Context, id string, hashedPassword string) error
	Delete(ctx context.Context, id string) error
}

type Service interface {
	Create(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, limit, offset int) ([]*User, int64, error)
	Update(ctx context.Context, id string, req *UpdateUserRequest) (*User, error)
	UpdateBalance(ctx context.Context, userID string, amount float64) error
	UpdatePassword(ctx context.Context, id string, hashedPassword string) error
	Delete(ctx context.Context, id string) error
}
