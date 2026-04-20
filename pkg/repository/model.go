package repository

import (
	"time"

	"github.com/zatrano/zatrano/pkg/tenant"
	"gorm.io/gorm"
)

// Model is the standard base model for all ZATRANO entities.
// Embed this in every model to get ID, timestamps, and soft-delete support.
//
//	type User struct {
//	    repository.Model
//	    Name  string
//	    Email string
//	}
type Model struct {
	ID        uint           `gorm:"primarykey"                  json:"id"`
	CreatedAt time.Time      `                                   json:"created_at"`
	UpdatedAt time.Time      `                                   json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"       json:"-"        swaggerignore:"true"`
}

// IsDeleted returns true if the record has been soft-deleted.
func (m *Model) IsDeleted() bool {
	return m.DeletedAt.Valid
}

// TenantFK is an optional embed for row-isolated models. When used with tenant.FromContext,
// BeforeCreate sets TenantID from the numeric tenant key (middleware.ResolveTenant + numeric X-Tenant-ID).
//
//	type Order struct {
//	    repository.Model
//	    repository.TenantFK
//	    Total int
//	}
type TenantFK struct {
	TenantID uint `gorm:"index;not null" json:"tenant_id"`
}

// BeforeCreate fills TenantID from context when the tenant key parses as a positive integer.
func (t *TenantFK) BeforeCreate(tx *gorm.DB) error {
	if tx.Statement == nil || tx.Statement.Context == nil {
		return nil
	}
	info, ok := tenant.FromContext(tx.Statement.Context)
	if !ok || info.NumericID == 0 {
		return nil
	}
	if t.TenantID == 0 {
		t.TenantID = uint(info.NumericID)
	}
	return nil
}
