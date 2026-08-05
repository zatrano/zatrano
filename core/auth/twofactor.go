package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/zatrano/framework/core/hashing"
	"github.com/zatrano/framework/core/http"
	"github.com/zatrano/framework/core/totp"
)

const twoFactorSessionKey = "auth.two_factor_user_id"

// ErrTwoFactorRequired indicates that the primary credentials are valid and a second factor is pending.
var ErrTwoFactorRequired = fmt.Errorf("two-factor authentication required")

type twoFactorUser interface {
	GetTwoFactorSecret() string
	GetTwoFactorRecoveryCodes() string
	GetTwoFactorConfirmedAt() *time.Time
}

func (m *Manager) twoFactorValues(user Authenticatable) (secret, recovery string, confirmed bool) {
	if user == nil {
		return "", "", false
	}
	var rawSecret, rawRecovery string
	if u, ok := user.(twoFactorUser); ok {
		rawSecret = u.GetTwoFactorSecret()
		rawRecovery = u.GetTwoFactorRecoveryCodes()
		confirmed = u.GetTwoFactorConfirmedAt() != nil
	} else if u, ok := user.(*GenericUser); ok {
		rawSecret = strings.TrimSpace(fmt.Sprint(u.Get("two_factor_secret")))
		if rawSecret == "<nil>" {
			rawSecret = ""
		}
		rawRecovery = strings.TrimSpace(fmt.Sprint(u.Get("two_factor_recovery_codes")))
		if rawRecovery == "<nil>" {
			rawRecovery = ""
		}
		confirmed = u.Get("two_factor_confirmed_at") != nil && fmt.Sprint(u.Get("two_factor_confirmed_at")) != "" && fmt.Sprint(u.Get("two_factor_confirmed_at")) != "<nil>"
	}
	return m.decryptSensitive(rawSecret), m.decryptSensitive(rawRecovery), confirmed
}

func (m *Manager) encryptSensitive(value string) string {
	if m == nil || m.crypt == nil || strings.TrimSpace(value) == "" {
		return value
	}
	out, err := m.crypt.Encrypt(value)
	if err != nil {
		return value
	}
	return out
}

func (m *Manager) decryptSensitive(value string) string {
	value = strings.TrimSpace(value)
	if m == nil || m.crypt == nil || value == "" {
		return value
	}
	out, err := m.crypt.Decrypt(value)
	if err != nil {
		// Allow plaintext values from before encryption was enabled.
		return value
	}
	return out
}

func (m *Manager) updateTwoFactor(user Authenticatable, attrs map[string]any) error {
	if m == nil || m.Guard() == nil {
		return fmt.Errorf("auth guard unavailable")
	}
	updater, ok := m.Guard().Provider().(AttributeUpdater)
	if !ok {
		return fmt.Errorf("user provider does not support two-factor updates")
	}
	if err := updater.UpdateAttributes(user.AuthID(), attrs); err != nil {
		return err
	}
	if generic, ok := user.(*GenericUser); ok {
		for key, value := range attrs {
			generic.Attributes[key] = value
		}
	}
	return nil
}

// EnableTwoFactor creates an unconfirmed TOTP secret and recovery codes.
func (m *Manager) EnableTwoFactor(user Authenticatable) (secret, otpauthURL string, recoveryCodes []string, err error) {
	if user == nil {
		return "", "", nil, fmt.Errorf("user is required")
	}
	secret, err = totp.GenerateSecret()
	if err != nil {
		return "", "", nil, err
	}
	recoveryCodes, err = makeRecoveryCodes()
	if err != nil {
		return "", "", nil, err
	}
	joined := strings.Join(recoveryCodes, ",")
	if err := m.updateTwoFactor(user, map[string]any{
		"two_factor_secret":         m.encryptSensitive(secret),
		"two_factor_recovery_codes": m.encryptSensitive(joined),
		"two_factor_confirmed_at":   nil,
	}); err != nil {
		return "", "", nil, err
	}
	otpauthURL = totp.OTPAuthURL("ZATRANO", EmailForVerification(user), secret)
	return secret, otpauthURL, recoveryCodes, nil
}

// ConfirmTwoFactor confirms a pending TOTP secret.
func (m *Manager) ConfirmTwoFactor(user Authenticatable, code string) error {
	secret, _, _ := m.twoFactorValues(user)
	if secret == "" || !totp.Verify(secret, strings.TrimSpace(code)) {
		return fmt.Errorf("invalid two-factor code")
	}
	return m.updateTwoFactor(user, map[string]any{"two_factor_confirmed_at": time.Now().UTC()})
}

// DisableTwoFactor removes all second-factor data after password confirmation.
func (m *Manager) DisableTwoFactor(user Authenticatable, password string) error {
	if user == nil || !hashing.Check(password, user.AuthPassword()) {
		return fmt.Errorf("current password is incorrect")
	}
	return m.updateTwoFactor(user, map[string]any{
		"two_factor_secret": nil, "two_factor_recovery_codes": nil, "two_factor_confirmed_at": nil,
	})
}

// HasTwoFactorEnabled reports whether a user has confirmed second-factor authentication.
func (m *Manager) HasTwoFactorEnabled(user Authenticatable) bool {
	secret, _, confirmed := m.twoFactorValues(user)
	return secret != "" && confirmed
}

// GenerateRecoveryCodes creates and stores a new recovery-code set.
func (m *Manager) GenerateRecoveryCodes(user Authenticatable) ([]string, error) {
	if !m.HasTwoFactorEnabled(user) {
		return nil, fmt.Errorf("two-factor authentication is not enabled")
	}
	codes, err := makeRecoveryCodes()
	if err != nil {
		return nil, err
	}
	return codes, m.updateTwoFactor(user, map[string]any{"two_factor_recovery_codes": m.encryptSensitive(strings.Join(codes, ","))})
}

// ReplaceRecoveryCodes is an alias for GenerateRecoveryCodes.
func (m *Manager) ReplaceRecoveryCodes(user Authenticatable) ([]string, error) {
	return m.GenerateRecoveryCodes(user)
}

// VerifyTwoFactorCode validates a TOTP code or consumes a recovery code.
func (m *Manager) VerifyTwoFactorCode(user Authenticatable, code string) bool {
	secret, recovery, confirmed := m.twoFactorValues(user)
	code = strings.TrimSpace(code)
	if !confirmed || secret == "" {
		return false
	}
	if totp.Verify(secret, code) {
		return true
	}
	codes := splitRecoveryCodes(recovery)
	for i, saved := range codes {
		if code != saved {
			continue
		}
		codes = append(codes[:i], codes[i+1:]...)
		return m.updateTwoFactor(user, map[string]any{"two_factor_recovery_codes": m.encryptSensitive(strings.Join(codes, ","))}) == nil
	}
	return false
}

// ChallengeTwoFactor completes a pending two-factor challenge and logs the user in.
func (m *Manager) ChallengeTwoFactor(req *http.Request, code string) (bool, error) {
	if req == nil || req.Session() == nil {
		return false, fmt.Errorf("session not available")
	}
	id := req.Session().Get(twoFactorSessionKey)
	if id == nil || fmt.Sprint(id) == "" {
		return false, fmt.Errorf("two-factor challenge is not pending")
	}
	user, err := m.Guard().Provider().RetrieveByID(id)
	if err != nil || user == nil {
		return false, err
	}
	if !m.VerifyTwoFactorCode(user, code) {
		return false, fmt.Errorf("invalid two-factor code")
	}
	req.Session().Forget(twoFactorSessionKey)
	if err := m.Login(req, user); err != nil {
		return false, err
	}
	m.dispatch(EventTwoFactorAuthenticated, TwoFactorAuthenticatedEvent{Request: req, User: user, Guard: m.Guard().name, At: time.Now().UTC()})
	return true, nil
}

func makeRecoveryCodes() ([]string, error) {
	codes := make([]string, 8)
	for i := range codes {
		raw := make([]byte, 5)
		if _, err := rand.Read(raw); err != nil {
			return nil, err
		}
		codes[i] = hex.EncodeToString(raw)
	}
	return codes, nil
}

func splitRecoveryCodes(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if code := strings.TrimSpace(part); code != "" {
			out = append(out, code)
		}
	}
	return out
}
