package database

import (
	"fmt"
	"strings"

	"github.com/zatrano/zatrano/internal/zatrano/query" // Query paketimiz

	"gorm.io/gorm"
)

// ApplyQuery, bir Query nesnesini bir GORM sorgusuna uygular.
func ApplyQuery(db *gorm.DB, q *query.Query) *gorm.DB {
	tx := db

	// Filtreleri uygula
	for _, filter := range q.Filters {
		// Operatörleri SQL'e çevir
		op := "=" // varsayılan
		val := filter.Value
		switch strings.ToLower(filter.Operator) {
		case "eq":
			op = "="
		case "neq":
			op = "!="
		case "gt":
			op = ">"
		case "gte":
			op = ">="
		case "lt":
			op = "<"
		case "lte":
			op = "<="
		case "like":
			op = "LIKE"
			val = fmt.Sprintf("%%%v%%", filter.Value)
		}
		tx = tx.Where(fmt.Sprintf("%s %s ?", filter.Field, op), val)
	}

	// Sıralamayı uygula
	for _, sort := range q.Sorts {
		tx = tx.Order(fmt.Sprintf("%s %s", sort.Field, sort.Direction))
	}

	// İlişkileri (Preload) uygula
	for _, relation := range q.Relations {
		tx = tx.Preload(relation)
	}

	return tx
}
