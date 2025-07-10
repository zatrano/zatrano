package models

import (
	"time"
	"gorm.io/gorm"
)

// BaseModel, tüm veritabanı modelleri için ortak alanları içerir.
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// --- YENİ DENETİM ALANLARI ---
	// gorm:"default:null" ile bu alanların boş olabilir (nullable) olduğunu belirtiyoruz.
	// Böylece sistemsel (otomatik) oluşturulan kayıtlarda bu alanlar boş kalabilir.
	CreatedBy *uint          `gorm:"default:null" json:"created_by,omitempty"`
	UpdatedBy *uint          `gorm:"default:null" json:"updated_by,omitempty"`
	DeletedBy *uint          `gorm:"default:null" json:"-"`
}