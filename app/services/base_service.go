package services

import (
	"github.com/jinzhu/copier"
	"github.com/zatrano/zatrano/app/repositories"
)

// IBaseService, jenerik CRUD işlemlerini form nesneleriyle birlikte yönetir.
// Bu arayüz, handler'ların kullanacağı sözleşmedir.
type IBaseService[T any, CForm any, UForm any] interface {
	GetAll(page, limit int) ([]T, int64, error)
	GetByID(id uint, relations ...string) (*T, error)
	CreateFromForm(form *CForm) (*T, error) // Form struct'ı alır, T modeli döndürür
	UpdateFromForm(id uint, form *UForm) (*T, error) // Form struct'ı alır, T modeli döndürür
	Delete(id uint) error
}

// baseService, arayüzün implementasyonudur.
type baseService[T any, CForm any, UForm any] struct {
	repo repositories.IBaseRepository[T]
}

func NewBaseService[T any, CForm any, UForm any](repo repositories.IBaseRepository[T]) IBaseService[T, CForm, UForm] {
	return &baseService[T, CForm, UForm]{repo: repo}
}

// ... GetAll, GetByID, Delete metodları önceki cevapla aynı ...
func (s *baseService[T, CForm, UForm]) GetAll(...) { /* ... */ }
func (s *baseService[T, CForm, UForm]) GetByID(...) { /* ... */ }
func (s *baseService[T, CForm, UForm]) Delete(...) { /* ... */ }

// CreateFromForm, bir form nesnesinden yeni bir kayıt oluşturur.
func (s *baseService[T, CForm, UForm]) CreateFromForm(form *CForm) (*T, error) {
	var model T
	// `copier` kullanarak form verilerini yeni ve boş bir modele kopyala.
	if err := copier.Copy(&model, form); err != nil {
		return nil, err // Kopyalama hatası
	}
	
	if err := s.repo.Create(&model); err != nil {
		return nil, err // Veritabanı hatası
	}
	return &model, nil
}

// UpdateFromForm, bir formu kullanarak mevcut bir kaydı günceller.
func (s *baseService[T, CForm, UForm]) UpdateFromForm(id uint, form *UForm) (*T, error) {
	// Önce güncellenecek kaydı veritabanından bul.
	model, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	
	// Gelen form verilerini, veritabanından gelen mevcut modelin üzerine kopyala.
	// `copier` varsayılan olarak boş ("zero-value") alanları atlar.
	// Yani formdaki şifre alanı boşsa, modeldeki şifreye dokunmaz.
	if err := copier.Copy(model, form); err != nil {
		return nil, err
	}

	// Artık güncellenmiş ve dolu olan modeli repository'ye gönder.
	if err := s.repo.Update(model); err != nil {
		return nil, err
	}
	return model, nil
}