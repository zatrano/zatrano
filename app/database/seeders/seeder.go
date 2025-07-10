package seeders

import "gorm.io/gorm"

// ISeeder, tüm seeder struct'larının uygulaması gereken arayüzdür.
type ISeeder interface {
	Run(db *gorm.DB) error
}