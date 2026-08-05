package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/zatrano/framework/core/http"
	"github.com/zatrano/framework/core/routing"
)

// MustVerifyEmail is implemented by users that require email verification.
type MustVerifyEmail interface {
	HasVerifiedEmail() bool
	GetEmailForVerification() string
	AuthID() any
}

// HasVerifiedEmail reports whether the authenticatable user has a verified email.
func HasVerifiedEmail(user Authenticatable) bool {
	if user == nil {
		return false
	}
	if v, ok := user.(MustVerifyEmail); ok {
		return v.HasVerifiedEmail()
	}
	if generic, ok := user.(*GenericUser); ok {
		raw := generic.Get("email_verified_at")
		if raw == nil {
			return false
		}
		s := strings.TrimSpace(fmt.Sprint(raw))
		return s != "" && s != "<nil>"
	}
	return true
}

// EmailForVerification returns the email used for verification links.
func EmailForVerification(user Authenticatable) string {
	if user == nil {
		return ""
	}
	if v, ok := user.(MustVerifyEmail); ok {
		return v.GetEmailForVerification()
	}
	if generic, ok := user.(*GenericUser); ok {
		return fmt.Sprint(generic.Get("email"))
	}
	return ""
}

// EmailHash returns a stable hash of the email for signed verification URLs.
func EmailHash(email string) string {
	sum := sha1.Sum([]byte(strings.ToLower(strings.TrimSpace(email))))
	return hex.EncodeToString(sum[:])
}

// VerifyEmailMiddleware blocks authenticated users who have not verified email.
func VerifyEmailMiddleware(manager *Manager) routing.MiddlewareFunc {
	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *http.Request) *http.Response {
			user := manager.User(req)
			if user == nil {
				if req.WantsJSON() {
					return http.JSON(map[string]any{"message": "Unauthenticated."}).Status(401)
				}
				return http.Redirect("/login")
			}
			if HasVerifiedEmail(user) {
				return next(req)
			}
			if req.WantsJSON() {
				return http.JSON(map[string]any{"message": "Your email address is not verified."}).Status(403)
			}
			return http.Redirect("/email/verify")
		}
	}
}

// MarkEmailVerified updates a GenericUser attribute map (caller persists).
func MarkEmailVerified(attrs map[string]any) {
	if attrs == nil {
		return
	}
	attrs["email_verified_at"] = time.Now().UTC()
}

// HasVerifiedEmail implements MustVerifyEmail for GenericUser.
func (u *GenericUser) HasVerifiedEmail() bool {
	if u == nil {
		return false
	}
	raw := u.Get("email_verified_at")
	if raw == nil {
		return false
	}
	s := strings.TrimSpace(fmt.Sprint(raw))
	return s != "" && s != "<nil>"
}

// GetEmailForVerification implements MustVerifyEmail for GenericUser.
func (u *GenericUser) GetEmailForVerification() string {
	if u == nil {
		return ""
	}
	return fmt.Sprint(u.Get("email"))
}
