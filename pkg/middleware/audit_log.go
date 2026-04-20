package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	jwtlib "github.com/golang-jwt/jwt/v5"

	"github.com/zatrano/zatrano/pkg/audit"
	"github.com/zatrano/zatrano/pkg/config"
	"github.com/zatrano/zatrano/pkg/security"
	"go.uber.org/zap"
)

// AuditLog records one line per HTTP request (database or JSONL file) when audit.http_enabled is on.
func AuditLog(cfg *config.Config, w audit.Writer, log *zap.Logger) fiber.Handler {
	if cfg == nil || w == nil || log == nil {
		log = zap.NewNop()
	}
	if !cfg.Audit.Enabled || !cfg.Audit.HttpEnabled {
		return func(c fiber.Ctx) error { return c.Next() }
	}
	return func(c fiber.Ctx) error {
		start := time.Now()
		ctx := c.Context()
		rid := requestid.FromContext(c)
		ctx = audit.WithRequest(ctx, rid, c.IP())
		if uid := fiberUserID(c); uid != "" {
			ctx = audit.WithUser(ctx, uid)
		} else if mc, ok := c.Locals(security.ClaimsKey()).(jwtlib.MapClaims); ok {
			if sub, ok := mc["sub"].(string); ok && sub != "" {
				ctx = audit.WithUser(ctx, sub)
			}
		}
		c.SetContext(ctx)

		err := c.Next()

		q := string(c.Request().URI().QueryString())
		var qPtr *string
		if q != "" {
			qPtr = &q
		}
		st := c.Response().StatusCode()
		dur := int(time.Since(start) / time.Millisecond)
		uid := audit.UserFromContext(ctx)
		var uidPtr *string
		if uid != "" {
			uidPtr = &uid
		}
		ridPtr := &rid
		if rid == "" {
			ridPtr = nil
		}
		ip := c.IP()
		ipPtr := &ip
		row := &audit.HTTPAuditLog{
			CreatedAt:  time.Now(),
			Method:     c.Method(),
			Path:       c.Path(),
			URLQuery:   qPtr,
			Status:     st,
			DurationMs: dur,
			UserID:     uidPtr,
			RequestID:  ridPtr,
			IP:         ipPtr,
		}
		if werr := w.WriteHTTP(ctx, row); werr != nil {
			log.Warn("audit http log", zap.Error(werr))
		}
		return err
	}
}

func fiberUserID(c fiber.Ctx) string {
	if v := c.Locals(LocalsUserID); v != nil {
		switch t := v.(type) {
		case uint:
			return fmt.Sprintf("%d", t)
		case int:
			return fmt.Sprintf("%d", t)
		case float64:
			return fmt.Sprintf("%.0f", t)
		case string:
			return t
		}
	}
	return ""
}
