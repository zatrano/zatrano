package query

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Filter, tek bir filtre koşulunu temsil eder.
// Örn: {Field: "name", Operator: "like", Value: "%john%"}
type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

// Sort, tek bir sıralama koşulunu temsil eder.
// Örn: {Field: "created_at", Direction: "desc"}
type Sort struct {
	Field     string
	Direction string
}

// Pagination, sayfalama verilerini tutar.
type Pagination struct {
	Page  int
	Limit int
}

// Query, bir HTTP isteğinden parse edilmiş tüm sorgu parametrelerini içerir.
type Query struct {
	Filters    []Filter
	Sorts      []Sort
	Pagination Pagination
	Relations  []string // Preload edilecek ilişkiler
}

// Parser, bir Fiber context'ini Query nesnesine dönüştürür.
type Parser struct {
	// Geliştiricinin hangi alanlara göre filtreleme/sıralama yapabileceğini
	// sınırlamak için kullanılır. Bu, güvenlik için önemlidir.
	AllowedFilters []string
	AllowedSorts   []string
}

func NewParser() *Parser {
	return &Parser{}
}

// Parse, bir Fiber context'ini analiz eder.
// URL Formatı:
// - Filtre: ?filter[name][like]=John&filter[status][eq]=active
// - Sıralama: ?sort=-created_at,name (eksi işareti DESC anlamına gelir)
// - Sayfalama: ?page[number]=2&page[size]=15
// - İlişkiler: ?include=user,comments
func (p *Parser) Parse(c *fiber.Ctx) *Query {
	q := &Query{
		Filters:    p.parseFilters(c),
		Sorts:      p.parseSorts(c),
		Pagination: p.parsePagination(c),
		Relations:  p.parseIncludes(c),
	}
	return q
}

func (p *Parser) parseFilters(c *fiber.Ctx) []Filter {
	var filters []Filter
	// `c.Queries()` "filter[name][like]" gibi anahtarları map'e çevirmez.
	// Bu yüzden `c.Request().URI().QueryArgs()` kullanmak daha esnektir.
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		keyStr := string(key)
		if strings.HasPrefix(keyStr, "filter[") && strings.HasSuffix(keyStr, "]") {
			// Örnek parse mantığı: filter[name][like]
			parts := strings.Split(strings.TrimSuffix(strings.TrimPrefix(keyStr, "filter["), "]"), "][")
			if len(parts) == 2 {
				field := parts[0]
				operator := parts[1]
				filters = append(filters, Filter{Field: field, Operator: operator, Value: string(value)})
			}
		}
	})
	return filters
}

func (p *Parser) parseSorts(c *fiber.Ctx) []Sort {
	var sorts []Sort
	sortQuery := c.Query("sort")
	if sortQuery == "" {
		return sorts
	}
	fields := strings.Split(sortQuery, ",")
	for _, field := range fields {
		direction := "asc"
		if strings.HasPrefix(field, "-") {
			direction = "desc"
			field = strings.TrimPrefix(field, "-")
		}
		sorts = append(sorts, Sort{Field: field, Direction: direction})
	}
	return sorts
}

func (p *Parser) parsePagination(c *fiber.Ctx) Pagination {
	page, _ := strconv.Atoi(c.Query("page[number]", "1"))
	limit, _ := strconv.Atoi(c.Query("page[size]", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	} // Maksimum limit
	return Pagination{Page: page, Limit: limit}
}

func (p *Parser) parseIncludes(c *fiber.Ctx) []string {
	includeQuery := c.Query("include")
	if includeQuery == "" {
		return nil
	}
	return strings.Split(includeQuery, ",")
}
