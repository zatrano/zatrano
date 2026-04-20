package features

import (
	"github.com/gofiber/fiber/v3"
)

// Locals keys — must match pkg/middleware/authorize.go (InjectRoles / RBAC).
const (
	fiberLocalsUserID    = "zatrano.auth.user_id"
	fiberLocalsUserRoles = "zatrano.auth.user_roles"
)

// UserFromFiber builds a features.User from RBAC Locals (may return nil user).
func UserFromFiber(c fiber.Ctx) *User {
	if c == nil {
		return nil
	}
	var uid uint
	if v := c.Locals(fiberLocalsUserID); v != nil {
		switch t := v.(type) {
		case uint:
			uid = t
		case int:
			if t > 0 {
				uid = uint(t)
			}
		case float64:
			if t > 0 {
				uid = uint(t)
			}
		}
	}
	var roles []string
	if v := c.Locals(fiberLocalsUserRoles); v != nil {
		if rr, ok := v.([]string); ok {
			roles = rr
		}
	}
	if uid == 0 && len(roles) == 0 {
		return nil
	}
	return &User{ID: uid, Roles: roles}
}

// FromFiber returns an Eval for the current request (anonymous when no Locals user).
func (r *Registry) FromFiber(c fiber.Ctx) *Eval {
	if r == nil {
		return &Eval{}
	}
	return r.For(UserFromFiber(c))
}

// Middleware attaches an *Eval to Fiber Locals for handlers that want it without templates.
const LocalsEval = "zatrano.features.eval"

// LocalsMiddleware stores *Eval under LocalsEval for the request.
func (r *Registry) LocalsMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		if r != nil && r.Enabled() {
			c.Locals(LocalsEval, r.FromFiber(c))
		}
		return c.Next()
	}
}

// EvalFromFiber returns the Eval from Locals (nil if middleware not used or features disabled).
func EvalFromFiber(c fiber.Ctx) *Eval {
	if c == nil {
		return nil
	}
	v := c.Locals(LocalsEval)
	if v == nil {
		return nil
	}
	e, _ := v.(*Eval)
	return e
}
