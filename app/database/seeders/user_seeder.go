package seeders

import (
	"fmt"
	"github.com/zatrano/zatrano/app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder, users tablosuna başlangıç verileri ekler.
type UserSeeder struct{}

func NewUserSeeder() ISeeder {
	return &UserSeeder{}
}

func (s *UserSeeder) Run(db *gorm.DB) error {
	// Şifreyi hash'le
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	// Oluşturulacak kullanıcılar
	users := []models.User{
		{
			BaseModel: models.BaseModel{ID: 1}, // ID'yi manuel olarak belirlemek için
			Name:      "Admin User",
			Email:     "admin@zatrano.com",
			Password:  string(password),
		},
		{
			BaseModel: models.BaseModel{ID: 2},
			Name:      "Test User",
			Email:     "test@zatrano.com",
			Password:  string(password),
		},
	}
	
	for _, user := range users {
		// Eğer bu ID'ye sahip bir kullanıcı zaten varsa, ekleme (idempotent)
		var count int64
		db.Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
		if count == 0 {
			err := db.Create(&user).Error
			if err != nil {
				return fmt.Errorf("error seeding user %s: %w", user.Email, err)
			}
		}
	}
	
	fmt.Println("User seeder ran successfully.")
	return nil
}