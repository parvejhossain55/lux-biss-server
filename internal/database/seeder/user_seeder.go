package seeder

import (
	"github.com/google/uuid"
	"github.com/parvej/luxbiss_server/internal/modules/user"
	"github.com/parvej/luxbiss_server/pkg/hash"
	"gorm.io/gorm"
)

type UserSeeder struct{}

func (s *UserSeeder) Seed(db *gorm.DB) error {
	hashedPassword, _ := hash.HashPassword("Parvej@55")

	users := []user.User{
		{
			ID:       uuid.New().String(),
			Name:     "Admin User",
			Email:    "admin@luxbiss.com",
			Password: hashedPassword,
			Role:     user.RoleAdmin,
			Status:   user.StatusActive,
		},
		{
			ID:       uuid.New().String(),
			Name:     "Regular User",
			Email:    "user@luxbiss.com",
			Password: hashedPassword,
			Role:     user.RoleUser,
			Status:   user.StatusActive,
		},
	}

	for _, u := range users {
		var existing user.User
		if err := db.Where("email = ?", u.Email).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&u).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
