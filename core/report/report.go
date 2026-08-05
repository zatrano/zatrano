package report

import (
	"sync"
	"time"

	"github.com/zatrano/framework/core/http"
)

// Event is a captured exception report.
type Event struct {
	ID        int64     `json:"id"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
	Path      string    `json:"path,omitempty"`
	Method    string    `json:"method,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Manager stores recent exception reports (Sentry-like stub).
type Manager struct {
	mu     sync.Mutex
	nextID int64
	events []Event
	limit  int
}

// New creates a report manager.
func New(limit ...int) *Manager {
	n := 100
	if len(limit) > 0 && limit[0] > 0 {
		n = limit[0]
	}
	return &Manager{limit: n, events: make([]Event, 0), nextID: 1}
}

// Capture records an error.
func (m *Manager) Capture(err error, req *http.Request, level ...string) Event {
	if err == nil {
		return Event{}
	}
	lvl := "error"
	if len(level) > 0 && level[0] != "" {
		lvl = level[0]
	}
	ev := Event{
		Message:   err.Error(),
		Level:     lvl,
		CreatedAt: time.Now().UTC(),
	}
	if req != nil {
		ev.Path = req.Path()
		ev.Method = req.Method()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	ev.ID = m.nextID
	m.nextID++
	m.events = append([]Event{ev}, m.events...)
	if len(m.events) > m.limit {
		m.events = m.events[:m.limit]
	}
	return ev
}

// Recent returns the latest events.
func (m *Manager) Recent(limit int) []Event {
	m.mu.Lock()
	defer m.mu.Unlock()
	if limit <= 0 || limit > len(m.events) {
		limit = len(m.events)
	}
	out := make([]Event, limit)
	copy(out, m.events[:limit])
	return out
}

// Count returns stored event count.
func (m *Manager) Count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.events)
}

// Clear removes all events.
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = m.events[:0]
}

// Reporter returns an exceptions.Reporter-compatible callback.
func (m *Manager) Reporter() func(err error, req *http.Request) {
	return func(err error, req *http.Request) {
		_ = m.Capture(err, req)
	}
}
