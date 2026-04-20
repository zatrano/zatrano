package tenant

import (
	"context"
	"strconv"
	"strings"
)

type ctxKey int

const infoKey ctxKey = iota + 1

// Info carries the resolved tenant for the request (HTTP + std context).
type Info struct {
	// Key is the raw tenant identifier from header or subdomain (trimmed).
	Key string
	// NumericID is parsed when Key is a base-10 unsigned integer (otherwise 0).
	NumericID uint64
	// Schema is the PostgreSQL schema name when isolation=schema (quoted-safe fragment).
	Schema string
}

// WithContext returns ctx storing tenant info for downstream GORM/repository use.
func WithContext(ctx context.Context, info Info) context.Context {
	return context.WithValue(ctx, infoKey, info)
}

// FromContext returns tenant info when middleware (or tests) attached it.
func FromContext(ctx context.Context) (Info, bool) {
	if ctx == nil {
		return Info{}, false
	}
	v, ok := ctx.Value(infoKey).(Info)
	if !ok || strings.TrimSpace(v.Key) == "" {
		return Info{}, false
	}
	return v, true
}

// ParseNumericKey sets NumericID when Key is a non-empty digit string.
func ParseNumericKey(key string) uint64 {
	key = strings.TrimSpace(key)
	if key == "" {
		return 0
	}
	n, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		return 0
	}
	return n
}
