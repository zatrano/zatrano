package audit

import (
	"context"
)

type ctxKey int

const (
	ctxKeyUserID ctxKey = iota + 1
	ctxKeyRequestID
	ctxKeyIP
	ctxKeySkip
)

// WithUser stores the acting user id (e.g. JWT sub) for model activity rows.
func WithUser(ctx context.Context, userID string) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return context.WithValue(ctx, ctxKeyUserID, userID)
}

// UserFromContext returns the user id set by WithUser or HTTP middleware.
func UserFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	s, _ := ctx.Value(ctxKeyUserID).(string)
	return s
}

// WithRequest attaches request correlation fields used by audit writers.
func WithRequest(ctx context.Context, requestID, ip string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = context.WithValue(ctx, ctxKeyRequestID, requestID)
	ctx = context.WithValue(ctx, ctxKeyIP, ip)
	return ctx
}

func requestFromContext(ctx context.Context) (requestID, ip string) {
	if ctx == nil {
		return "", ""
	}
	requestID, _ = ctx.Value(ctxKeyRequestID).(string)
	ip, _ = ctx.Value(ctxKeyIP).(string)
	return requestID, ip
}

// SkipNext disables audit callbacks for the current GORM session chain when set on context.
func Skip(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return context.WithValue(ctx, ctxKeySkip, true)
}

func skipFromContext(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	v, _ := ctx.Value(ctxKeySkip).(bool)
	return v
}
