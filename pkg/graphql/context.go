package graphql

import (
	"context"
)

type loadersCtxKey struct{}

// WithLoaders attaches per-request loaders to ctx (call from the GraphQL HTTP wrapper).
func WithLoaders(ctx context.Context, l *Loaders) context.Context {
	if l == nil {
		return ctx
	}
	return context.WithValue(ctx, loadersCtxKey{}, l)
}

// LoadersFrom returns loaders previously attached with WithLoaders, or nil.
func LoadersFrom(ctx context.Context) *Loaders {
	if ctx == nil {
		return nil
	}
	v, _ := ctx.Value(loadersCtxKey{}).(*Loaders)
	return v
}
