package repositories

import (
	"gorm.io/gorm"
)

// IBaseRepository, herhangi bir model T için temel veritabanı işlemlerini tanımlar.
type IBaseRepository[T any] interface {
	FindAll(page, limit int) ([]T, int64, error)
	FindByID(id uint, relations ...string) (*T, error)
	Create(model *T) error
	Update(model *T) error
	Delete(model *T) error
}

// baseRepository, IBaseRepository arayüzünün GORM implementasyonudur.
type baseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository, jenerik repository için bir yapıcı fonksiyondur.
func NewBaseRepository[T any](db *gorm.DB) IBaseRepository[T] {
	return &baseRepository[T]{db: db}
}

// FindAll, tüm kayıtları paginated olarak alır.
func (r *baseRepository[T]) FindAll(page, limit int) ([]T, int64, error) {
	var models []T
	var total int64
	offset := (page - 1) * limit

	modelType := new(T)
	if err := r.db.Model(modelType).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Limit(limit).Offset(offset).Find(&models).Error
	return models, total, err
}

// FindByID, tek bir kaydı ID'ye ve istenen ilişkilere göre bulur.
func (r *baseRepository[T]) FindByID(id uint, relations ...string) (*T, error) {
	var model T
	query := r.db
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	err := query.First(&model, id).Error
	return &model, err
}

// Create, yeni bir kayıt oluşturur.
func (r *baseRepository[T]) Create(model *T) error {
	return r.db.Create(model).Error
}

// Update, mevcut bir kaydın tüm alanlarını günceller.
// GORM, sadece değişen alanları güncelleyecektir.
// GORM hook'larının çalışması için dolu model gönderilir.
func (r *baseRepository[T]) Update(model *T) error {
	return r.db.Save(model).Error
}

// Delete, bir kaydı soft-delete olarak işaretler.
// GORM hook'larının çalışması için dolu model gönderilir.
func (r *baseRepository[T]) Delete(model *T) error {
	return r.db.Delete(model).Error
}