package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// BruteForceProtection handles login attempt rate limiting
type BruteForceProtection struct {
	redis       *redis.Client
	maxAttempts int
	window      time.Duration
	blockTime   time.Duration
}

// NewBruteForceProtection creates a new brute force protection instance
func NewBruteForceProtection(redisClient *redis.Client) *BruteForceProtection {
	return &BruteForceProtection{
		redis:       redisClient,
		maxAttempts: 5,                // 5 failed attempts
		window:      15 * time.Minute, // within 15 minutes
		blockTime:   30 * time.Minute, // block for 30 minutes
	}
}

// CheckLoginAttempt checks if a login attempt should be allowed
func (b *BruteForceProtection) CheckLoginAttempt(identifier string) error {
	ctx := context.Background()
	key := fmt.Sprintf("login_attempts:%s", identifier)

	// Get current attempt count
	count, err := b.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to check login attempts: %w", err)
	}

	// Check if user is currently blocked
	blockKey := fmt.Sprintf("login_blocked:%s", identifier)
	blocked, err := b.redis.Exists(ctx, blockKey).Result()
	if err != nil {
		return fmt.Errorf("failed to check block status: %w", err)
	}

	if blocked > 0 {
		return fmt.Errorf("account temporarily blocked due to too many failed login attempts")
	}

	// Check if max attempts exceeded
	if count >= b.maxAttempts {
		// Block the user
		if err := b.redis.Set(ctx, blockKey, "1", b.blockTime).Err(); err != nil {
			return fmt.Errorf("failed to block user: %w", err)
		}
		return fmt.Errorf("account temporarily blocked due to too many failed login attempts")
	}

	return nil
}

// RecordFailedAttempt records a failed login attempt
func (b *BruteForceProtection) RecordFailedAttempt(identifier string) error {
	ctx := context.Background()
	key := fmt.Sprintf("login_attempts:%s", identifier)

	// Increment attempt count
	count, err := b.redis.Incr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to record failed attempt: %w", err)
	}

	// Set expiry on the key if this is the first attempt
	if count == 1 {
		if err := b.redis.Expire(ctx, key, b.window).Err(); err != nil {
			return fmt.Errorf("failed to set expiry on attempt counter: %w", err)
		}
	}

	return nil
}

// RecordSuccessfulLogin clears failed attempts for a successful login
func (b *BruteForceProtection) RecordSuccessfulLogin(identifier string) error {
	ctx := context.Background()
	key := fmt.Sprintf("login_attempts:%s", identifier)
	blockKey := fmt.Sprintf("login_blocked:%s", identifier)

	// Remove attempt counter and block
	pipe := b.redis.Pipeline()
	pipe.Del(ctx, key)
	pipe.Del(ctx, blockKey)
	_, err := pipe.Exec(ctx)

	return err
}

// GetRemainingAttempts returns the number of remaining login attempts
func (b *BruteForceProtection) GetRemainingAttempts(identifier string) (int, error) {
	ctx := context.Background()
	key := fmt.Sprintf("login_attempts:%s", identifier)

	count, err := b.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return 0, fmt.Errorf("failed to get attempt count: %w", err)
	}

	remaining := b.maxAttempts - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// IsBlocked checks if an identifier is currently blocked
func (b *BruteForceProtection) IsBlocked(identifier string) (bool, error) {
	ctx := context.Background()
	blockKey := fmt.Sprintf("login_blocked:%s", identifier)

	blocked, err := b.redis.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check block status: %w", err)
	}

	return blocked > 0, nil
}
