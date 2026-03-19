package product

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Level struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" gorm:"type:varchar(50);not null;unique"`
	ProfitPercentage float64        `json:"profit_percentage" gorm:"type:decimal(5,2);not null;default:0"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Level) TableName() string {
	return "levels"
}

type Step struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	LevelID    uint           `json:"level_id" gorm:"not null;index"`
	Level      *Level         `json:"level,omitempty" gorm:"foreignKey:LevelID"`
	StepNumber int            `json:"step_number" gorm:"not null"`
	Name       string         `json:"name" gorm:"type:varchar(50)"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Step) TableName() string {
	return "steps"
}

type Product struct {
	ID          string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	LevelID     uint           `json:"level_id" gorm:"not null;index"`
	Level       *Level         `json:"level,omitempty" gorm:"foreignKey:LevelID"`
	StepID      uint           `json:"step_id" gorm:"not null;index"`
	Step        *Step          `json:"step,omitempty" gorm:"foreignKey:StepID"`
	Name        string         `json:"name" gorm:"type:varchar(255);not null"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	Rating      float64        `json:"rating" gorm:"type:decimal(2,1);default:0.0"`
	MinQuantity int            `json:"min_quantity" gorm:"not null;default:1"`
	MaxQuantity int            `json:"max_quantity" gorm:"not null;default:100"`
	ImageURL    string         `json:"image_url" gorm:"type:text"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}

type Repository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	List(ctx context.Context, limit, offset int, sortBy, order string, levelID, stepID uint) ([]*Product, int64, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error

	// Level operations
	CreateLevel(ctx context.Context, level *Level) error
	GetLevelByID(ctx context.Context, id uint) (*Level, error)
	ListLevels(ctx context.Context, limit, offset int) ([]*Level, int64, error)
	UpdateLevel(ctx context.Context, level *Level) error
	DeleteLevel(ctx context.Context, id uint) error

	// Step operations
	CreateStep(ctx context.Context, step *Step) error
	GetStepByID(ctx context.Context, id uint) (*Step, error)
	ListStepsByLevel(ctx context.Context, levelID uint, limit, offset int) ([]*Step, int64, error)
	UpdateStep(ctx context.Context, step *Step) error
	DeleteStep(ctx context.Context, id uint) error
}

type Service interface {
	Create(ctx context.Context, req *CreateProductRequest) (*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	List(ctx context.Context, limit, offset int, sortBy, order string, levelID, stepID uint) ([]*Product, int64, error)
	Update(ctx context.Context, id string, req *UpdateProductRequest) (*Product, error)
	Delete(ctx context.Context, id string) error

	// Level operations
	CreateLevel(ctx context.Context, req *CreateLevelRequest) (*Level, error)
	GetLevelByID(ctx context.Context, id uint) (*Level, error)
	ListLevels(ctx context.Context, limit, offset int) ([]*Level, int64, error)
	UpdateLevel(ctx context.Context, id uint, req *UpdateLevelRequest) (*Level, error)
	DeleteLevel(ctx context.Context, id uint) error

	// Step operations
	CreateStep(ctx context.Context, req *CreateStepRequest) (*Step, error)
	GetStepByID(ctx context.Context, id uint) (*Step, error)
	ListStepsByLevel(ctx context.Context, levelID uint, limit, offset int) ([]*Step, int64, error)
	UpdateStep(ctx context.Context, id uint, req *UpdateStepRequest) (*Step, error)
	DeleteStep(ctx context.Context, id uint) error
}
