package seeder

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/parvej/luxbiss_server/internal/modules/product"
	"gorm.io/gorm"
)

type ProductSeeder struct{}

func (s *ProductSeeder) Seed(db *gorm.DB) error {
	var levels []product.Level
	if err := db.Find(&levels).Error; err != nil {
		return err
	}

	for _, level := range levels {
		var steps []product.Step
		if err := db.Where("level_id = ?", level.ID).Order("step_number ASC").Find(&steps).Error; err != nil {
			return err
		}

		for _, step := range steps {
			numProducts := 1
			// Steps 4, 8, 12, 16, 20 are "combo" steps with 4 products
			if step.StepNumber%4 == 0 {
				numProducts = 4
			}

			for i := 1; i <= numProducts; i++ {
				productName := fmt.Sprintf("%s %s Item %d", level.Name, step.Name, i)
				if numProducts == 4 {
					productName = fmt.Sprintf("%s %s Combo %d", level.Name, step.Name, i)
				}

				var existing product.Product
				if err := db.Where("name = ? AND level_id = ? AND step_id = ?", productName, level.ID, step.ID).First(&existing).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						p := product.Product{
							ID:          uuid.New().String(),
							LevelID:     level.ID,
							StepID:      step.ID,
							Name:        productName,
							Price:       float64(50 + (level.ID * 5) + uint(step.StepNumber*2) + uint(i*10)),
							Rating:      4.0 + (float64(i) * 0.2),
							MinQuantity: 1,
							MaxQuantity: 50,
							Description: fmt.Sprintf("Premium %s specifically curated for %s and %s.", productName, level.Name, step.Name),
							ImageURL:    fmt.Sprintf("https://picsum.photos/seed/prod-%d-%d-%d/600/400", level.ID, step.ID, i),
						}
						if err := db.Create(&p).Error; err != nil {
							return err
						}
					} else {
						return err
					}
				}
			}
		}
	}

	return nil
}
