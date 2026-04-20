package broadcast

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/zatrano/zatrano/pkg/config"
)

// ParseAccessToken validates HS256 JWT the same way as security.JWTMiddleware.
func ParseAccessToken(cfg *config.Config, raw string) (jwt.MapClaims, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("missing token")
	}
	secret := []byte(strings.TrimSpace(cfg.Security.JWTSecret))
	if len(secret) == 0 {
		return nil, fmt.Errorf("jwt_secret not configured")
	}
	issuer := cfg.Security.JWTIssuer
	tok, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	}, jwt.WithIssuer(issuer), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !tok.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	mc, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return mc, nil
}

// Subject returns JWT "sub" as string (may be empty).
func Subject(mc jwt.MapClaims) string {
	if mc == nil {
		return ""
	}
	v, ok := mc["sub"]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case float64:
		return fmt.Sprintf("%.0f", t)
	default:
		return fmt.Sprint(t)
	}
}
