package repository

import (
	"context"

	"github.com/zatrano/zatrano/pkg/tenant"
	"gorm.io/gorm"
)

// TenantScope returns a Scope that restricts queries to the current tenant (row isolation).
// column is the SQL column name (e.g. tenant_id or tenant_slug). Numeric tenant keys use
// uint comparison when ParseNumericKey succeeds; otherwise the string Key is used.
func TenantScope(ctx context.Context, column string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		info, ok := tenant.FromContext(ctx)
		if !ok || column == "" {
			return db
		}
		if info.NumericID > 0 {
			return db.Where(column+" = ?", uint(info.NumericID))
		}
		if info.Key != "" {
			return db.Where(column+" = ?", info.Key)
		}
		return db
	}
}
