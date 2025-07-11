package services

import (
	"github.com/jinzhu/copier"
	"github.com/zatrano/zatrano/app/repositories"
	"github.com/zatrano/zatrano/internal/zatrano/query"
)

// IBaseService, basit ve gelişmiş CRUD işlemlerini, form nesneleriyle birlikte yönetir.
type IBaseService[T any, CForm any, UForm any] interface {
	// === MEVCUT METODLAR (Korunuyor ve Geliştiriliyor) ===
	GetAll(page, limit int) ([]T, int64, error)
	GetByID(id uint, relations ...string) (*T, error)
	GetByCondition(condition map[string]interface{}, relations ...string) (*T, error)
	FindAllByCondition(condition map[string]interface{}, page, limit int, orderBy string, relations ...string) ([]T, int64, error)
	CreateFromForm(form *CForm) (*T, error)
	UpdateFromForm(id uint, form *UForm) (*T, error)
	Delete(id uint) error
	BulkDelete(condition map[string]interface{}) error

	// === YENİ EKLENEN MERKEZİ METODLAR ===
	Find(q *query.Query) ([]T, int64, error)
	FindOne(q *query.Query) (*T, error)
}

// baseService, IBaseService arayüzünün implementasyonudur.
type baseService[T any, CForm any, UForm any] struct {
	repo repositories.IBaseRepository[T]
}

func NewBaseService[T any, CForm any, UForm any](repo repositories.IBaseRepository[T]) IBaseService[T, CForm, UForm] {
	return &baseService[T, CForm, UForm]{repo: repo}
}

// ===================================================================
// OKUMA İŞLEMLERİ (Mevcut Metodlar)
// ===================================================================

func (s *baseService[T, CForm, UForm]) GetAll(page, limit int) ([]T, int64, error) {
	return s.repo.GetAll(page, limit)
}

func (s *baseService[T, CForm, UForm]) GetByID(id uint, relations ...string) (*T, error) {
	return s.repo.GetByID(id, relations...)
}

func (s *baseService[T, CForm, UForm]) GetByCondition(condition map[string]interface{}, relations ...string) (*T, error) {
	return s.repo.GetByCondition(condition, relations...)
}

func (s *baseService[T, CForm, UForm]) FindAllByCondition(condition map[string]interface{}, page, limit int, orderBy string, relations ...string) ([]T, int64, error) {
	return s.repo.FindAllByCondition(condition, page, limit, orderBy, relations...)
}

// ===================================================================
// OKUMA İŞLEMLERİ (Yeni Merkezi Metodlar)
// ===================================================================

// Find, merkezi Query nesnesini alır ve repository'ye iletir.
// Servis katmanı, bu sorguya ek iş kuralları (örn: yetkilendirme) ekleyebilir.
func (s *baseService[T, CForm, UForm]) Find(q *query.Query) ([]T, int64, error) {
	// Örnek iş kuralı: Eğer giriş yapmış bir kullanıcı varsa ve bu bir multi-tenant
	// uygulama ise, sorguya otomatik olarak tenant ID filtresi ekle.
	// if authUser := auth.User(); authUser != nil {
	// 	q.Filters = append(q.Filters, query.Filter{
	// 		Field: "tenant_id", Operator: "eq", Value: authUser.TenantID,
	// 	})
	// }
	return s.repo.Find(q)
}

// FindOne, merkezi Query nesnesini kullanarak tek bir kayıt bulur.
func (s *baseService[T, CForm, UForm]) FindOne(q *query.Query) (*T, error) {
	return s.repo.FindOne(q)
}

// ===================================================================
// YAZMA İŞLEMLERİ
// ===================================================================

func (s *baseService[T, CForm, UForm]) CreateFromForm(form *CForm) (*T, error) {
	var model T
	if err := copier.Copy(&model, form); err != nil {
		return nil, err
	}
	if err := s.repo.Create(&model); err != nil {
		return nil, err
	}
	// event.Dispatch(new ModelCreatedEvent(model))
	return &model, nil
}

func (s *baseService[T, CForm, UForm]) UpdateFromForm(id uint, form *UForm) (*T, error) {
	model, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if err := copier.Copy(model, form); err != nil {
		return nil, err
	}
	if err := s.repo.Update(model); err != nil {
		return nil, err
	}
	// cache.Forget("model_key:" + id)
	return model, nil
}

// ===================================================================
// SİLME İŞLEMLERİ
// ===================================================================

func (s *baseService[T, CForm, UForm]) Delete(id uint) error {
	model, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	// if auth.User().Cant("delete", model) { return errors.New("unauthorized") }
	return s.repo.Delete(model)
}

// ===================================================================
// TOPLU (BULK) İŞLEMLER
// ===================================================================

func (s *baseService[T, CForm, UForm]) BulkDelete(condition map[string]interface{}) error {
	// if !auth.User().IsAdmin() { return errors.New("only admins can perform bulk delete") }
	return s.repo.BulkDelete(condition)
}
