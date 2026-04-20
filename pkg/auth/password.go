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

// PasswordResetService handles password reset functionality
type PasswordResetService struct {
	db   *gorm.DB
	mail *mail.Manager
	cfg  *config.Config
}

// NewPasswordResetService creates a new password reset service
func NewPasswordResetService(db *gorm.DB, mailer *mail.Manager, cfg *config.Config) *PasswordResetService {
	return &PasswordResetService{
		db:   db,
		mail: mailer,
		cfg:  cfg,
	}
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        uint      `gorm:"primarykey"`
	Email     string    `gorm:"not null;index"`
	Token     string    `gorm:"not null;uniqueIndex;size:255"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}

// GenerateResetToken generates a secure reset token and stores it
func (s *PasswordResetService) GenerateResetToken(email string) (string, error) {
	// Generate secure token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Store token
	resetToken := &PasswordResetToken{
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
	}

	if err := s.db.Create(resetToken).Error; err != nil {
		return "", fmt.Errorf("failed to store reset token: %w", err)
	}

	return token, nil
}

// ValidateResetToken validates a reset token
func (s *PasswordResetService) ValidateResetToken(token string) (string, error) {
	var resetToken PasswordResetToken
	if err := s.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&resetToken).Error; err != nil {
		return "", fmt.Errorf("invalid or expired token")
	}

	return resetToken.Email, nil
}

// ResetPassword resets the password for the given email
func (s *PasswordResetService) ResetPassword(token, newPassword string) error {
	email, err := s.ValidateResetToken(token)
	if err != nil {
		return err
	}

	_ = email // TODO: Use email to update user password

	// Here you would update the user's password
	// This depends on your user model structure
	// For example:
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	// if err != nil {
	//     return fmt.Errorf("failed to hash password: %w", err)
	// }
	// if err := s.db.Model(&User{}).Where("email = ?", email).Update("password", string(hashedPassword)).Error; err != nil {
	//     return fmt.Errorf("failed to update password: %w", err)
	// }

	// Delete used token
	if err := s.db.Where("token = ?", token).Delete(&PasswordResetToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete used token: %w", err)
	}

	return nil
}

// SendResetEmail sends password reset email
func (s *PasswordResetService) SendResetEmail(email, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", "https://"+s.cfg.HTTPAddr, token)

	subject := "Password Reset Request"
	body := fmt.Sprintf(`
Hello,

You have requested to reset your password. Click the link below to reset your password:

%s

This link will expire in 1 hour.

If you did not request this password reset, please ignore this email.

Best regards,
%s Team
`, resetURL, s.cfg.AppName)

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

// CleanExpiredTokens removes expired reset tokens
func (s *PasswordResetService) CleanExpiredTokens() error {
	return s.db.Where("expires_at < ?", time.Now()).Delete(&PasswordResetToken{}).Error
}
