package auth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"net/url"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/pquerna/otp/totp"
	"github.com/zatrano/zatrano/pkg/config"
	"gorm.io/gorm"
)

// TOTPService handles TOTP (Time-based One-Time Password) 2FA
type TOTPService struct {
	cfg *config.Config
	db  *gorm.DB
}

// NewTOTPService creates a new TOTP service
func NewTOTPService(cfg *config.Config, db *gorm.DB) *TOTPService {
	return &TOTPService{
		cfg: cfg,
		db:  db,
	}
}

// TOTPSetup represents TOTP setup data for a user
type TOTPSetup struct {
	Secret     string `json:"secret"`
	QRCodeURL  string `json:"qr_code_url"`
	QRCodeData string `json:"qr_code_data"` // Base64 encoded PNG
}

// GenerateTOTPSecret generates a new TOTP secret for a user
func (s *TOTPService) GenerateTOTPSecret(userID uint, email string) (*TOTPSetup, error) {
	// Generate TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.cfg.AppName,
		AccountName: email,
		SecretSize:  32,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// Generate QR code
	qrCode, err := s.generateQRCode(key.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	setup := &TOTPSetup{
		Secret:     key.Secret(),
		QRCodeURL:  key.URL(),
		QRCodeData: qrCode,
	}

	// Store secret temporarily (user needs to verify before enabling)
	// In a real implementation, you might store this in a temporary table
	// or cache until verification is complete

	return setup, nil
}

// VerifyTOTPCode verifies a TOTP code
func (s *TOTPService) VerifyTOTPCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

// EnableTOTP enables TOTP for a user after successful verification
func (s *TOTPService) EnableTOTP(userID uint, secret string) error {
	// Update user with TOTP secret
	// This assumes your User model has a totp_secret field
	if err := s.db.Model(&User{}).Where("id = ?", userID).Update("totp_secret", secret).Error; err != nil {
		return fmt.Errorf("failed to enable TOTP: %w", err)
	}

	return nil
}

// DisableTOTP disables TOTP for a user
func (s *TOTPService) DisableTOTP(userID uint) error {
	// Clear TOTP secret
	if err := s.db.Model(&User{}).Where("id = ?", userID).Update("totp_secret", nil).Error; err != nil {
		return fmt.Errorf("failed to disable TOTP: %w", err)
	}

	return nil
}

// IsTOTPEnabled checks if TOTP is enabled for a user
func (s *TOTPService) IsTOTPEnabled(userID uint) (bool, error) {
	var user User
	if err := s.db.Select("totp_secret").Where("id = ?", userID).First(&user).Error; err != nil {
		return false, err
	}

	return user.TOTPSecret != nil && *user.TOTPSecret != "", nil
}

// ValidateUserTOTP validates TOTP for a user during login
func (s *TOTPService) ValidateUserTOTP(userID uint, code string) (bool, error) {
	var user User
	if err := s.db.Select("totp_secret").Where("id = ?", userID).First(&user).Error; err != nil {
		return false, err
	}

	if user.TOTPSecret == nil || *user.TOTPSecret == "" {
		// TOTP not enabled for this user
		return true, nil
	}

	return s.VerifyTOTPCode(*user.TOTPSecret, code), nil
}

// generateQRCode generates a base64 encoded QR code PNG
func (s *TOTPService) generateQRCode(url string) (string, error) {
	// Create QR code
	qrCode, err := qr.Encode(url, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}

	// Scale the barcode to 256x256
	qrCode, err = barcode.Scale(qrCode, 256, 256)
	if err != nil {
		return "", err
	}

	// Encode to PNG
	var buf bytes.Buffer
	err = png.Encode(&buf, qrCode)
	if err != nil {
		return "", err
	}

	// Convert to base64
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// GetTOTPURL generates a TOTP URL for manual entry
func (s *TOTPService) GetTOTPURL(email, secret string) string {
	u := &url.URL{
		Scheme: "otpauth",
		Host:   "totp",
		Path:   fmt.Sprintf("/%s:%s", s.cfg.AppName, email),
	}
	q := u.Query()
	q.Set("secret", secret)
	q.Set("issuer", s.cfg.AppName)
	u.RawQuery = q.Encode()

	return u.String()
}
