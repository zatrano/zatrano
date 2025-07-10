package seeders

import (
	"fmt"
	"gorm.io/gorm"
)

// DatabaseSeeder, tüm seeder'ları sırayla çalıştırır.
type DatabaseSeeder struct {
	seeders []ISeeder
}

// NewDatabaseSeeder, çalıştırılacak tüm seeder'ları alır.
func NewDatabaseSeeder() *DatabaseSeeder {
	return &DatabaseSeeder{
		seeders: []ISeeder{
			// Bağımlılıkları olanlar daha sonra çalışmalı.
			// Örn: Bir post'un bir user'ı olmalı, o yüzden önce UserSeeder.
			NewUserSeeder(),
			// NewCountrySeeder(),
			// NewProductSeeder(),
		},
	}
}

// Run, tüm seeder'ları sırayla çalıştırır.
func (s *DatabaseSeeder) Run(db *gorm.DB) error {
	for _, seeder := range s.seeders {
		if err := seeder.Run(db); err != nil {
			return fmt.Errorf("seeding failed: %w", err)
		}
	}
	fmt.Println("All seeders ran successfully!")
	return nil
}