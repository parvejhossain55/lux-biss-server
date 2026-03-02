package seeder

import (
	"fmt"

	"github.com/parvej/luxbiss_server/internal/modules/product"
	"gorm.io/gorm"
)

type LevelSeeder struct{}

func (s *LevelSeeder) Seed(db *gorm.DB) error {
	for i := 1; i <= 20; i++ {
		name := fmt.Sprintf("Level %d", i)
		var existing product.Level
		if err := db.Where("name = ?", name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				level := product.Level{
					Name: name,
				}
				if err := db.Create(&level).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}
