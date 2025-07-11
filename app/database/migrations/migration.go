package migrations

import (
	"log"

	"github.com/zatrano/zatrano/app/models"
	"gorm.io/gorm"
)

// RunMigrations, tüm veritabanı migration'larını çalıştırır.
// GORM'un AutoMigrate özelliği bu iş için mükemmeldir.
func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// GORM, bu modellerin karşılığı olan tabloları oluşturur veya
	// mevcut tablolara eksik sütunları/indeksleri ekler.
	// Tablo veya sütun silmez.
	err := db.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.Country{},
		&models.Session{},
	)

	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Migrations completed successfully.")
	return nil
}
