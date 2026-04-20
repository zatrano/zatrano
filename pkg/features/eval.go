package features

import (
	"context"
	"hash/fnv"
	"strconv"
	"strings"
)

// Eval evaluates flags for a fixed user binding.
type Eval struct {
	reg  *Registry
	user *User
}

// IsEnabled resolves the flag for this evaluator's user (see also Registry.Resolve).
func (e *Eval) IsEnabled(ctx context.Context, key string) bool {
	if e == nil || e.reg == nil {
		return false
	}
	return e.reg.resolveEnabled(ctx, e.user, key)
}

func (r *Registry) resolveEnabled(ctx context.Context, u *User, key string) bool {
	def, ok := r.resolveDefinition(ctx, key)
	if !ok {
		return false
	}
	if !def.Enabled {
		return false
	}
	if len(def.AllowedRoles) > 0 {
		if u == nil {
			return false
		}
		if !rolesIntersect(u.Roles, def.AllowedRoles) {
			return false
		}
	}
	if def.RolloutPercent <= 0 {
		return true
	}
	if def.RolloutPercent >= 100 {
		return true
	}
	if u == nil || u.ID == 0 {
		// Partial rollout requires a stable user id for A/B buckets.
		return false
	}
	return rolloutAllows(key, u.ID, def.RolloutPercent)
}

func rolesIntersect(userRoles, allowed []string) bool {
	seen := make(map[string]struct{}, len(userRoles))
	for _, r := range userRoles {
		seen[strings.ToLower(strings.TrimSpace(r))] = struct{}{}
	}
	for _, a := range allowed {
		a = strings.ToLower(strings.TrimSpace(a))
		if a == "" {
			continue
		}
		if _, ok := seen[a]; ok {
			return true
		}
	}
	return false
}

func rolloutAllows(key string, userID uint, pct int) bool {
	h := fnv.New32a()
	_, _ = h.Write([]byte(normalizeKey(key)))
	_, _ = h.Write([]byte(strconv.FormatUint(uint64(userID), 10)))
	b := h.Sum32() % 100
	return int(b) < pct
}

// Resolve returns the effective on/off for a user without constructing Eval (same logic).
func (r *Registry) Resolve(ctx context.Context, u *User, key string) bool {
	if r == nil {
		return false
	}
	return r.resolveEnabled(ctx, u, key)
}
