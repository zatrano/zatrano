package features

import "strings"

// Definition is the resolved flag shape used for evaluation.
type Definition struct {
	Enabled        bool
	RolloutPercent int
	AllowedRoles   []string
}

func normalizeKey(key string) string {
	return strings.ToLower(strings.TrimSpace(key))
}
