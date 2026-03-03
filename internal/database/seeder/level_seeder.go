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
				existing = level
			} else {
				return err
			}
		}

		// Create 20 steps for each level
		for j := 1; j <= 20; j++ {
			var existingStep product.Step
			if err := db.Where("level_id = ? AND step_number = ?", existing.ID, j).First(&existingStep).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					step := product.Step{
						LevelID:    existing.ID,
						StepNumber: j,
						Name:       fmt.Sprintf("Step %d", j),
					}
					if err := db.Create(&step).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}
		}
	}
	return nil
}
