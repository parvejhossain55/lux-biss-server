package product

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Level struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(50);not null;unique"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Level) TableName() string {
	return "levels"
}

type Product struct {
	ID          string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	LevelID     uint           `json:"level_id" gorm:"not null;index"`
	Level       *Level         `json:"level,omitempty" gorm:"foreignKey:LevelID"`
	Name        string         `json:"name" gorm:"type:varchar(255);not null"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	Rating      float64        `json:"rating" gorm:"type:decimal(2,1);default:0.0"`
	MinQuantity int            `json:"min_quantity" gorm:"not null;default:1"`
	MaxQuantity int            `json:"max_quantity" gorm:"not null;default:100"`
	ImageURL    string         `json:"image_url" gorm:"type:varchar(255)"`
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
	List(ctx context.Context, limit, offset int) ([]*Product, int64, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
}

type Service interface {
	Create(ctx context.Context, req *CreateProductRequest) (*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	List(ctx context.Context, limit, offset int) ([]*Product, int64, error)
	Update(ctx context.Context, id string, req *UpdateProductRequest) (*Product, error)
	Delete(ctx context.Context, id string) error
}
