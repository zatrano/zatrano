package services

import "github.com/zatrano/zatrano/app/repositories"

// IBaseService, herhangi bir model T için temel iş mantığı işlemlerini tanımlar.
type IBaseService[T any] interface {
	GetAll(page, limit int) ([]T, int64, error)
	GetByID(id uint, relations ...string) (*T, error)
	Create(model *T) error
	Update(id uint, updatedModel *T) (*T, error)
	Delete(id uint) error
}

type baseService[T any] struct {
	repo repositories.IBaseRepository[T]
}

// NewBaseService, jenerik servis için bir yapıcı fonksiyondur.
func NewBaseService[T any](repo repositories.IBaseRepository[T]) IBaseService[T] {
	return &baseService[T]{repo: repo}
}

// GetAll, tüm kayıtları repository'den alır.
func (s *baseService[T]) GetAll(page, limit int) ([]T, int64, error) {
	return s.repo.FindAll(page, limit)
}

// GetByID, tek bir kaydı repository'den alır.
func (s *baseService[T]) GetByID(id uint, relations ...string) (*T, error) {
	return s.repo.FindByID(id, relations...)
}

// Create, yeni bir kaydı repository aracılığıyla oluşturur.
func (s *baseService[T]) Create(model *T) error {
	// Gelecekte burada event fırlatma, loglama, cache temizleme gibi işlemler olabilir.
	return s.repo.Create(model)
}

// Update, bir kaydı günceller.
func (s *baseService[T]) Update(id uint, updatedModel *T) (*T, error) {
	// Önce güncellenecek kaydın varlığını kontrol et.
	existingModel, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err // Kayıt bulunamadı.
	}

	// Gelen `updatedModel`'deki ID'nin, URL'den gelen ID ile aynı olduğunu garantile.
	// Bu, GORM'un yanlış kaydı güncellemesini önler.
	// Reflection kullanarak ID'yi ayarlamak gerekir.
	setIDField(updatedModel, existingModel)

	if err := s.repo.Update(updatedModel); err != nil {
		return nil, err
	}
	return updatedModel, nil
}

// Delete, bir kaydı siler.
func (s *baseService[T]) Delete(id uint) error {
	// Önce silinecek kaydın varlığını kontrol et.
	modelToDelete, err := s.repo.FindByID(id)
	if err != nil {
		return err // Kayıt bulunamadı.
	}
	return s.repo.Delete(modelToDelete)
}

// setIDField, reflection kullanarak bir modelin ID alanını ayarlar.
// Bu, `Update` işleminde güvenliği sağlamak için gereklidir.
import "reflect"

func setIDField(target, source interface{}) {
	targetVal := reflect.ValueOf(target).Elem()
	sourceVal := reflect.ValueOf(source).Elem()

	idField := sourceVal.FieldByName("ID")
	if idField.IsValid() {
		targetIDField := targetVal.FieldByName("ID")
		if targetIDField.IsValid() && targetIDField.CanSet() {
			targetIDField.Set(idField)
		} else {
            // Embed edilmiş BaseModel için
            baseModelField := targetVal.FieldByName("BaseModel")
            if baseModelField.IsValid() {
                targetIDField = baseModelField.FieldByName("ID")
                if targetIDField.IsValid() && targetIDField.CanSet() {
                    targetIDField.Set(idField)
                }
            }
        }
	}
}