package search

import (
	"strings"

	"github.com/zatrano/zatrano/pkg/repository"
	"gorm.io/gorm"
)

// WhereFullTextMatch adds `vectorCol @@ plainto_tsquery(regconfig, query)` when query is non-empty.
// vectorCol must be a trusted column identifier (e.g. "search_vector"); it is interpolated as raw SQL.
func WhereFullTextMatch(regconfig, vectorCol, query string) repository.Scope {
	reg := strings.TrimSpace(regconfig)
	if reg == "" {
		reg = "simple"
	}
	col := strings.TrimSpace(vectorCol)
	return func(db *gorm.DB) *gorm.DB {
		q := strings.TrimSpace(query)
		if q == "" || col == "" {
			return db
		}
		return db.Where(col+" @@ plainto_tsquery(?::regconfig, ?)", reg, q)
	}
}

// OrderByTSRank orders by ts_rank_cd when query is non-empty.
func OrderByTSRank(regconfig, vectorCol, query string) repository.Scope {
	reg := strings.TrimSpace(regconfig)
	if reg == "" {
		reg = "simple"
	}
	col := strings.TrimSpace(vectorCol)
	return func(db *gorm.DB) *gorm.DB {
		q := strings.TrimSpace(query)
		if q == "" || col == "" {
			return db
		}
		return db.Order(gorm.Expr("ts_rank_cd("+col+", plainto_tsquery(?::regconfig, ?)) DESC", reg, q))
	}
}
