package tenant

import (
	"fmt"
	"regexp"
	"strings"
)

var safeKey = regexp.MustCompile(`^[a-z0-9][a-z0-9_]{0,62}$`)

// SanitizeKey normalizes a tenant slug for use inside PostgreSQL schema names.
// Returns lowercased [a-z0-9_] or error if empty/invalid after normalization.
func SanitizeKey(raw string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(raw))
	s = strings.ReplaceAll(s, "-", "_")
	if s == "" {
		return "", fmt.Errorf("empty tenant key")
	}
	if !safeKey.MatchString(s) {
		return "", fmt.Errorf("tenant key must match [a-z0-9_]+ (got %q)", raw)
	}
	return s, nil
}

// SchemaName returns prefixed schema identifier for isolation=schema.
func SchemaName(prefix, key string) (string, error) {
	k, err := SanitizeKey(key)
	if err != nil {
		return "", err
	}
	p := strings.TrimSpace(prefix)
	if p == "" {
		p = "tenant_"
	}
	return p + k, nil
}
