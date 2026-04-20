package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zatrano/zatrano/pkg/config"
	"gorm.io/gorm"
)

// JWTService handles JWT token operations with refresh token support
type JWTService struct {
	cfg *config.Config
	db  *gorm.DB
}

// NewJWTService creates a new JWT service
func NewJWTService(cfg *config.Config, db *gorm.DB) *JWTService {
	return &JWTService{
		cfg: cfg,
		db:  db,
	}
}

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"not null;uniqueIndex;size:255" json:"token"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// GenerateTokenPair generates both access and refresh tokens
func (s *JWTService) GenerateTokenPair(userID uint, extraClaims map[string]interface{}) (*TokenPair, error) {
	// Generate access token
	accessToken, err := s.generateAccessToken(userID, extraClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.cfg.Security.JWTExpiry.Seconds()),
	}, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *JWTService) RefreshAccessToken(refreshTokenStr string) (*TokenPair, error) {
	// Validate refresh token
	var refreshToken RefreshToken
	if err := s.db.Where("token = ? AND expires_at > ? AND revoked = false", refreshTokenStr, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// Generate new token pair
	tokenPair, err := s.GenerateTokenPair(refreshToken.UserID, nil)
	if err != nil {
		return nil, err
	}

	// Revoke old refresh token and create new one
	if err := s.revokeRefreshToken(refreshTokenStr); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	return tokenPair, nil
}

// RevokeRefreshToken revokes a refresh token
func (s *JWTService) RevokeRefreshToken(token string) error {
	return s.revokeRefreshToken(token)
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (s *JWTService) RevokeAllUserTokens(userID uint) error {
	return s.db.Where("user_id = ?", userID).Update("revoked", true).Error
}

// CleanExpiredTokens removes expired refresh tokens
func (s *JWTService) CleanExpiredTokens() error {
	return s.db.Where("expires_at < ?", time.Now()).Delete(&RefreshToken{}).Error
}

// generateAccessToken generates a JWT access token
func (s *JWTService) generateAccessToken(userID uint, extraClaims map[string]interface{}) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  fmt.Sprintf("%d", userID),
		"iat":  now.Unix(),
		"exp":  now.Add(s.cfg.Security.JWTExpiry).Unix(),
		"iss":  s.cfg.Security.JWTIssuer,
		"type": "access",
	}

	// Add extra claims
	for k, v := range extraClaims {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Security.JWTSecret))
}

// generateRefreshToken generates and stores a refresh token
func (s *JWTService) generateRefreshToken(userID uint) (string, error) {
	// Generate token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Store in database
	refreshToken := &RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	if err := s.db.Create(refreshToken).Error; err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return token, nil
}

// revokeRefreshToken marks a refresh token as revoked
func (s *JWTService) revokeRefreshToken(token string) error {
	return s.db.Model(&RefreshToken{}).Where("token = ?", token).Update("revoked", true).Error
}
