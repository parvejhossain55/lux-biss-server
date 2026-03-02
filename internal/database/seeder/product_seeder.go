package seeder

import (
	"github.com/google/uuid"
	"github.com/parvej/luxbiss_server/internal/modules/product"
	"gorm.io/gorm"
)

type ProductSeeder struct{}

func (s *ProductSeeder) Seed(db *gorm.DB) error {
	products := []product.Product{
		{
			ID:          uuid.New().String(),
			Name:        "Classic Leather Watch",
			Price:       199.99,
			Rating:      4.5,
			MinQuantity: 1,
			MaxQuantity: 10,
			ImageURL:    "https://example.com/images/watch.jpg",
			Description: "A timeless classic leather watch for every occasion.",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Premium Wireless Headphones",
			Price:       299.99,
			Rating:      4.8,
			MinQuantity: 1,
			MaxQuantity: 5,
			ImageURL:    "https://example.com/images/headphones.jpg",
			Description: "Experience crystal clear sound with our premium wireless headphones.",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Ergonomic Office Chair",
			Price:       449.50,
			Rating:      4.2,
			MinQuantity: 1,
			MaxQuantity: 20,
			ImageURL:    "https://example.com/images/chair.jpg",
			Description: "Work in comfort with our top-of-the-line ergonomic office chair.",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Smart Home Hub",
			Price:       129.00,
			Rating:      4.0,
			MinQuantity: 1,
			MaxQuantity: 15,
			ImageURL:    "https://example.com/images/hub.jpg",
			Description: "Control your entire home from one central smart hub.",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Professional Drone",
			Price:       899.99,
			Rating:      4.9,
			MinQuantity: 1,
			MaxQuantity: 2,
			ImageURL:    "https://example.com/images/drone.jpg",
			Description: "Capture breathtaking 4K footage with our professional-grade drone.",
		},
	}

	for _, p := range products {
		var existing product.Product
		if err := db.Where("name = ?", p.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&p).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
