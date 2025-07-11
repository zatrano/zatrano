package repositories

import (
	"github.com/zatrano/zatrano/internal/zatrano/database"
	"github.com/zatrano/zatrano/internal/zatrano/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IBaseRepository, basit, gelişmiş ve merkezi sorgulama metodlarını bir arada sunar.
type IBaseRepository[T any] interface {
	// === Sizin Mevcut Metodlarınız (Korunuyor) ===
	FindAll(page, limit int) ([]T, int64, error) // Page ve total'ı birleştirelim
	GetByID(id uint, relations ...string) (*T, error)
	GetByCondition(condition map[string]interface{}, relations ...string) (*T, error)
	GetCount() (int64, error)
	GetCountByCondition(condition map[string]interface{}) (int64, error)
	GetByConditionWithOrder(condition map[string]interface{}, orderBy string, relations ...string) (*T, error)
	FindAllByCondition(condition map[string]interface{}, page, limit int, orderBy string, relations ...string) ([]T, int64, error)
	Create(model *T) error
	CreateWithRelations(model *T) error
	Update(model *T) error
	BulkCreate(models []*T, batchSize int) error
	BulkUpdate(condition map[string]interface{}, data map[string]interface{}) error
	Delete(model *T) error
	DeleteWithRelations(model *T, relations ...string) error
	BulkDelete(condition map[string]interface{}) error

	// === YENİ EKLENEN MERKEZİ METODLAR ===
	Find(q *query.Query) ([]T, int64, error)
	FindOne(q *query.Query) (*T, error)
}

// baseRepository, IBaseRepository arayüzünün GORM implementasyonudur.
type baseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository, jenerik repository için bir yapıcı fonksiyondur.
func NewBaseRepository[T any](db *gorm.DB) IBaseRepository[T] {
	return &baseRepository[T]{db: db}
}

// ===================================================================
// MEVCUT METODLARIN IMPLEMENTASYONU (Sizin Kodunuzdan)
// ===================================================================

func (r *baseRepository[T]) FindAll(page, limit int) ([]T, int64, error) {
	var models []T
	var total int64
	offset := (page - 1) * limit
	r.db.Model(new(T)).Count(&total)
	err := r.db.Limit(limit).Offset(offset).Find(&models).Error
	return models, total, err
}

func (r *baseRepository[T]) FindAllByCondition(condition map[string]interface{}, page, limit int, orderBy string, relations ...string) ([]T, int64, error) {
	var models []T
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(new(T)).Where(condition)
	for _, relation := range relations {
		query = query.Preload(relation)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	err := query.Limit(limit).Offset(offset).Find(&models).Error
	return models, total, err
}

func (r *baseRepository[T]) GetByID(id uint, relations ...string) (*T, error) {
	var model T
	query := r.db
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	err := query.First(&model, id).Error
	return &model, err
}

func (r *baseRepository[T]) GetByCondition(condition map[string]interface{}, relations ...string) (*T, error) {
	var model T
	query := r.db.Where(condition)
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	err := query.First(&model).Error
	return &model, err
}

func (r *baseRepository[T]) GetByConditionWithOrder(condition map[string]interface{}, orderBy string, relations ...string) (*T, error) {
	var model T
	query := r.db.Where(condition)
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	if orderBy != "" {
		query = query.Order(orderBy)
	}
	err := query.First(&model).Error
	return &model, err
}

func (r *baseRepository[T]) GetCount() (int64, error) {
	var total int64
	err := r.db.Model(new(T)).Count(&total).Error
	return total, err
}

func (r *baseRepository[T]) GetCountByCondition(condition map[string]interface{}) (int64, error) {
	var total int64
	err := r.db.Model(new(T)).Where(condition).Count(&total).Error
	return total, err
}

func (r *baseRepository[T]) Create(model *T) error {
	return r.db.Create(model).Error
}

func (r *baseRepository[T]) CreateWithRelations(model *T) error {
	return r.db.Clauses(clause.Associations).Create(model).Error
}

func (r *baseRepository[T]) Update(model *T) error {
	return r.db.Save(model).Error
}

func (r *baseRepository[T]) BulkCreate(models []*T, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 1000
	}
	return r.db.CreateInBatches(models, batchSize).Error
}

func (r *baseRepository[T]) BulkUpdate(condition map[string]interface{}, data map[string]interface{}) error {
	var model T
	return r.db.Model(&model).Where(condition).Updates(data).Error
}

func (r *baseRepository[T]) Delete(model *T) error {
	return r.db.Delete(model).Error
}

func (r *baseRepository[T]) DeleteWithRelations(model *T, relations ...string) error {
	return r.db.Select(relations).Delete(model).Error
}

func (r *baseRepository[T]) BulkDelete(condition map[string]interface{}) error {
	var model T
	return r.db.Where(condition).Delete(&model).Error
}

// ===================================================================
// YENİ EKLENEN MERKEZİ QUERY METODLARI
// ===================================================================

// Find, merkezi Query nesnesini kullanarak verileri bulur, sayar ve sayfalar.
func (r *baseRepository[T]) Find(q *query.Query) ([]T, int64, error) {
	var models []T
	var total int64

	queryBuilder := r.db.Model(new(T))

	// Merkezi query parser'dan gelen koşulları, sıralamayı ve ilişkileri uygula
	queryBuilder = database.ApplyQuery(queryBuilder, q)

	// Sayfalama öncesi toplam kayıt sayısını al
	if err := queryBuilder.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sayfalamayı (limit ve offset) uygula
	offset := (q.Pagination.Page - 1) * q.Pagination.Limit
	queryBuilder = queryBuilder.Limit(q.Pagination.Limit).Offset(offset)

	// Sonuçları veritabanından çek
	if err := queryBuilder.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return models, total, nil
}

// FindOne, merkezi Query nesnesini kullanarak tek bir kayıt bulur.
func (r *baseRepository[T]) FindOne(q *query.Query) (*T, error) {
	var model T
	queryBuilder := r.db.Model(new(T))
	queryBuilder = database.ApplyQuery(queryBuilder, q)

	if err := queryBuilder.First(&model).Error; err != nil {
		return nil, err
	}

	return &model, nil
}
