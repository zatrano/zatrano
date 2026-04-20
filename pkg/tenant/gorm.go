package tenant

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// GormSession returns db scoped with PostgreSQL search_path when Info.Schema is set.
// Otherwise returns db.WithContext(ctx).
func GormSession(db *gorm.DB, ctx context.Context) *gorm.DB {
	if db == nil {
		return nil
	}
	if ctx == nil {
		return db
	}
	info, ok := FromContext(ctx)
	if !ok || strings.TrimSpace(info.Schema) == "" {
		return db.WithContext(ctx)
	}
	sdb := db.Session(&gorm.Session{PrepareStmt: true}).WithContext(ctx)
	q := fmt.Sprintf(`SET LOCAL search_path TO %q, public`, info.Schema)
	if err := sdb.Exec(q).Error; err != nil {
		return sdb
	}
	return sdb
}
