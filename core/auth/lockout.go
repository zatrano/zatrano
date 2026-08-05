package auth

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/zatrano/framework/core/http"
)

// ErrLockout is returned when login attempts are temporarily throttled.
var ErrLockout = fmt.Errorf("too many login attempts")

type lockoutRecord struct {
	attempts  int
	expiresAt time.Time
}

type lockoutStore struct {
	mu    sync.Mutex
	items map[string]lockoutRecord
	max   int
	decay time.Duration
}

func newLockoutStore(max int, decay time.Duration) *lockoutStore {
	if max <= 0 {
		max = 5
	}
	if decay <= 0 {
		decay = time.Minute
	}
	return &lockoutStore{items: make(map[string]lockoutRecord), max: max, decay: decay}
}

func (s *lockoutStore) locked(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[key]
	if !ok || time.Now().After(item.expiresAt) {
		delete(s.items, key)
		return false
	}
	return item.attempts >= s.max
}

func (s *lockoutStore) hit(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	item := s.items[key]
	if time.Now().After(item.expiresAt) {
		item = lockoutRecord{}
	}
	item.attempts++
	item.expiresAt = time.Now().Add(s.decay)
	s.items[key] = item
	return item.attempts >= s.max
}

func (s *lockoutStore) clear(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
}

func lockoutKey(req *http.Request, credentials map[string]string) string {
	ip := ""
	if req != nil {
		ip = req.IP()
	}
	return strings.ToLower(strings.TrimSpace(credentials["email"])) + "|" + ip
}
