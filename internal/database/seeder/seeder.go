package seeder

import "gorm.io/gorm"

type Seeder interface {
	Seed(db *gorm.DB) error
}

type Registry struct {
	seeders []Seeder
}

func NewRegistry() *Registry {
	return &Registry{
		seeders: make([]Seeder, 0),
	}
}

func (r *Registry) Register(s Seeder) {
	r.seeders = append(r.seeders, s)
}

func (r *Registry) RunAll(db *gorm.DB) error {
	for _, s := range r.seeders {
		if err := s.Seed(db); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) TruncateAll(db *gorm.DB, tables []string) error {
	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			return err
		}
	}
	return nil
}
