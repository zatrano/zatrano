package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/zatrano/zatrano/pkg/config"
	"github.com/zatrano/zatrano/pkg/mail"
	"gorm.io/gorm"
)

// EmailVerificationService handles email verification functionality
type EmailVerificationService struct {
	db   *gorm.DB
	mail *mail.Manager
	cfg  *config.Config
}

// NewEmailVerificationService creates a new email verification service
func NewEmailVerificationService(db *gorm.DB, mailer *mail.Manager, cfg *config.Config) *EmailVerificationService {
	return &EmailVerificationService{
		db:   db,
		mail: mailer,
		cfg:  cfg,
	}
}

// GenerateVerificationToken generates a secure verification token
func (s *EmailVerificationService) GenerateVerificationToken(userID uint, email string) (string, error) {
	// Generate secure token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Store token in user table (assuming you have email_verification_token and email_verification_expires_at columns)
	expiresAt := time.Now().Add(24 * time.Hour) // 24 hours expiry

	// Update user with verification token
	// This assumes your User model has these fields
	if err := s.db.Model(&User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"email_verification_token":      token,
		"email_verification_expires_at": expiresAt,
	}).Error; err != nil {
		return "", fmt.Errorf("failed to store verification token: %w", err)
	}

	return token, nil
}

// VerifyEmail verifies the email with the given token
func (s *EmailVerificationService) VerifyEmail(token string) error {
	// Find user with this token
	var user User
	if err := s.db.Where("email_verification_token = ? AND email_verification_expires_at > ?", token, time.Now()).First(&user).Error; err != nil {
		return fmt.Errorf("invalid or expired verification token")
	}

	// Mark email as verified
	now := time.Now()
	if err := s.db.Model(&user).Updates(map[string]interface{}{
		"email_verified_at":             &now,
		"email_verification_token":      nil,
		"email_verification_expires_at": nil,
	}).Error; err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

// SendVerificationEmail sends email verification email
func (s *EmailVerificationService) SendVerificationEmail(email, token string) error {
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", "https://"+s.cfg.HTTPAddr, token)

	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`
Hello,

Thank you for registering! Please verify your email address by clicking the link below:

%s

This link will expire in 24 hours.

If you did not create this account, please ignore this email.

Best regards,
%s Team
`, verificationURL, s.cfg.AppName)

	msg := &mail.Message{
		From: mail.Address{
			Name:  s.cfg.AppName,
			Email: s.cfg.Mail.FromEmail,
		},
		To: []mail.Address{
			{Email: email},
		},
		Subject:  subject,
		HTMLBody: body,
	}

	return s.mail.Send(context.Background(), msg)
}

// IsEmailVerified checks if the user's email is verified
func (s *EmailVerificationService) IsEmailVerified(userID uint) (bool, error) {
	var user User
	if err := s.db.Select("email_verified_at").Where("id = ?", userID).First(&user).Error; err != nil {
		return false, err
	}

	return user.EmailVerifiedAt != nil, nil
}

// User represents a basic user model (adjust according to your actual model)
type User struct {
	ID                         uint
	Email                      string
	EmailVerifiedAt            *time.Time
	EmailVerificationToken     *string
	EmailVerificationExpiresAt *time.Time
	TOTPSecret                 *string
}
