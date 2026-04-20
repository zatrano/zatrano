package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/zatrano/zatrano/pkg/config"
	"gorm.io/gorm"
)

// APIKey represents an API key for external client authentication.
type APIKey struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	Name      string     `gorm:"size:100;not null" json:"name"`
	Key       string     `gorm:"size:64;uniqueIndex;not null" json:"-"` // hashed
	Prefix    string     `gorm:"size:8;not null" json:"prefix"`
	Scopes    []string   `gorm:"serializer:json" json:"scopes"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName specifies the database table name.
func (APIKey) TableName() string { return "api_keys" }

// IsExpired checks if the key has expired.
func (k *APIKey) IsExpired() bool {
	return k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt)
}

// HasScope checks if the key has a specific scope.
func (k *APIKey) HasScope(scope string) bool {
	for _, s := range k.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// GenerateKey creates a new API key with prefix and hash.
func GenerateKey(name string, scopes []string, expiresAt *time.Time) (*APIKey, string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return nil, "", fmt.Errorf("generate random: %w", err)
	}
	plain := hex.EncodeToString(raw)
	hash := sha256.Sum256([]byte(plain))
	prefix := plain[:8]

	key := &APIKey{
		Name:      name,
		Key:       hex.EncodeToString(hash[:]),
		Prefix:    prefix,
		Scopes:    scopes,
		ExpiresAt: expiresAt,
	}
	return key, plain, nil
}

// ValidateKey checks if a provided key matches the stored hash.
func ValidateKey(storedHash, providedKey string) bool {
	hash := sha256.Sum256([]byte(providedKey))
	return storedHash == hex.EncodeToString(hash[:])
}

// Manager handles API key operations.
type KeyManager struct {
	db *gorm.DB
}

// NewKeyManager creates a new API key manager.
func NewKeyManager(db *gorm.DB) *KeyManager {
	return &KeyManager{db: db}
}

// Create generates and stores a new API key.
func (m *KeyManager) Create(name string, scopes []string, expiresAt *time.Time) (*APIKey, string, error) {
	key, plain, err := GenerateKey(name, scopes, expiresAt)
	if err != nil {
		return nil, "", err
	}
	if err := m.db.Create(key).Error; err != nil {
		return nil, "", fmt.Errorf("create api key: %w", err)
	}
	return key, plain, nil
}

// Authenticate validates an API key from request header.
func (m *KeyManager) Authenticate(headerValue string) (*APIKey, error) {
	if headerValue == "" {
		return nil, errors.New("missing api key")
	}
	parts := strings.SplitN(headerValue, ".", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid api key format")
	}
	prefix, providedKey := parts[0], parts[1]

	var key APIKey
	if err := m.db.Where("prefix = ? AND expires_at IS NULL OR expires_at > ?", prefix, time.Now()).First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid api key")
		}
		return nil, fmt.Errorf("lookup api key: %w", err)
	}

	if !ValidateKey(key.Key, providedKey) {
		return nil, errors.New("invalid api key")
	}

	return &key, nil
}

// Revoke marks an API key as expired immediately.
func (m *KeyManager) Revoke(id uint) error {
	return m.db.Model(&APIKey{}).Where("id = ?", id).Update("expires_at", time.Now()).Error
}

// List retrieves all active API keys.
func (m *KeyManager) List() ([]APIKey, error) {
	var keys []APIKey
	err := m.db.Where("expires_at IS NULL OR expires_at > ?", time.Now()).Find(&keys).Error
	return keys, err
}

// Middleware returns Fiber middleware for API key authentication.
func (m *KeyManager) Middleware(cfg *config.Config) fiber.Handler {
	header := cfg.Security.APIKeyHeader
	return func(c fiber.Ctx) error {
		if !cfg.Security.APIKeysEnabled {
			return c.Next()
		}
		key, err := m.Authenticate(c.Get(header))
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		c.Locals("api_key", key)
		return c.Next()
	}
}
