package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/zatrano/zatrano/pkg/config"
	"gorm.io/gorm"
)

// SessionService handles user session management
type SessionService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewSessionService creates a new session service
func NewSessionService(db *gorm.DB, cfg *config.Config) *SessionService {
	return &SessionService{
		db:  db,
		cfg: cfg,
	}
}

// UserSession represents a user session
type UserSession struct {
	ID           uint                   `gorm:"primarykey" json:"id"`
	UserID       uint                   `gorm:"not null;index" json:"user_id"`
	SessionToken string                 `gorm:"not null;uniqueIndex;size:255" json:"session_token"`
	RefreshToken *string                `gorm:"uniqueIndex;size:255" json:"refresh_token"`
	IPAddress    string                 `gorm:"type:inet" json:"ip_address"`
	UserAgent    string                 `gorm:"type:text" json:"user_agent"`
	DeviceInfo   map[string]interface{} `gorm:"type:jsonb" json:"device_info"`
	ExpiresAt    time.Time              `gorm:"not null;index" json:"expires_at"`
	LastActivity time.Time              `gorm:"default:current_timestamp" json:"last_activity"`
	CreatedAt    time.Time              `json:"created_at"`
}

// CreateSession creates a new user session
func (s *SessionService) CreateSession(userID uint, c fiber.Ctx) (*UserSession, error) {
	// Generate session token
	sessionToken, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Get client info
	ip := c.IP()
	userAgent := c.Get("User-Agent")

	session := &UserSession{
		UserID:       userID,
		SessionToken: sessionToken,
		RefreshToken: &refreshToken,
		IPAddress:    ip,
		UserAgent:    userAgent,
		DeviceInfo:   s.extractDeviceInfo(c),
		ExpiresAt:    time.Now().Add(24 * time.Hour), // 24 hours
		LastActivity: time.Now(),
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// ValidateSession validates a session token and updates last activity
func (s *SessionService) ValidateSession(sessionToken string) (*UserSession, error) {
	var session UserSession
	if err := s.db.Where("session_token = ? AND expires_at > ?", sessionToken, time.Now()).First(&session).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired session")
	}

	// Update last activity
	if err := s.db.Model(&session).Update("last_activity", time.Now()).Error; err != nil {
		return nil, fmt.Errorf("failed to update last activity: %w", err)
	}

	return &session, nil
}

// RefreshSession refreshes a session using refresh token
func (s *SessionService) RefreshSession(refreshToken string) (*UserSession, error) {
	var session UserSession
	if err := s.db.Where("refresh_token = ? AND expires_at > ?", refreshToken, time.Now()).First(&session).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// Generate new tokens
	newSessionToken, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new session token: %w", err)
	}

	newRefreshToken, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	// Update session
	session.SessionToken = newSessionToken
	session.RefreshToken = &newRefreshToken
	session.ExpiresAt = time.Now().Add(24 * time.Hour)
	session.LastActivity = time.Now()

	if err := s.db.Save(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	return &session, nil
}

// RevokeSession revokes a specific session
func (s *SessionService) RevokeSession(sessionToken string) error {
	result := s.db.Where("session_token = ?", sessionToken).Delete(&UserSession{})
	if result.Error != nil {
		return fmt.Errorf("failed to revoke session: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// RevokeAllUserSessions revokes all sessions for a user
func (s *SessionService) RevokeAllUserSessions(userID uint) error {
	return s.db.Where("user_id = ?", userID).Delete(&UserSession{}).Error
}

// GetUserSessions returns all active sessions for a user
func (s *SessionService) GetUserSessions(userID uint) ([]UserSession, error) {
	var sessions []UserSession
	if err := s.db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	return sessions, nil
}

// CleanExpiredSessions removes expired sessions
func (s *SessionService) CleanExpiredSessions() error {
	return s.db.Where("expires_at < ?", time.Now()).Delete(&UserSession{}).Error
}

// generateToken generates a secure random token
func (s *SessionService) generateToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
}

// extractDeviceInfo extracts device information from request headers
func (s *SessionService) extractDeviceInfo(c fiber.Ctx) map[string]interface{} {
	return map[string]interface{}{
		"user_agent":      c.Get("User-Agent"),
		"accept_language": c.Get("Accept-Language"),
		"accept_encoding": c.Get("Accept-Encoding"),
		"platform":        c.Get("Sec-Ch-Ua-Platform"),
		"mobile":          c.Get("Sec-Ch-Ua-Mobile"),
	}
}
