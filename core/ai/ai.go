package ai

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/zatrano/framework/core/support/uuid"
)

// Message is a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is a chat completion request.
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse is a chat completion response.
type ChatResponse struct {
	ID      string    `json:"id"`
	Model   string    `json:"model"`
	Message Message   `json:"message"`
	Usage   Usage     `json:"usage"`
	Created time.Time `json:"created_at"`
}

// Usage tracks token usage (stub counts).
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Driver generates completions.
type Driver interface {
	Name() string
	Chat(req ChatRequest) (*ChatResponse, error)
}

// Manager resolves AI drivers.
type Manager struct {
	mu            sync.RWMutex
	defaultDriver string
	drivers       map[string]Driver
}

// New creates an AI manager with a fake driver.
func New() *Manager {
	m := &Manager{drivers: make(map[string]Driver), defaultDriver: "fake"}
	m.Extend("fake", FakeDriver{})
	m.Extend("log", LogDriver{})
	return m
}

// Extend registers a driver.
func (m *Manager) Extend(name string, driver Driver) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.drivers[strings.ToLower(name)] = driver
}

// Use sets the default driver.
func (m *Manager) Use(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.defaultDriver = strings.ToLower(name)
}

// Driver returns a named driver.
func (m *Manager) Driver(name ...string) (Driver, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n := m.defaultDriver
	if len(name) > 0 && name[0] != "" {
		n = strings.ToLower(name[0])
	}
	d, ok := m.drivers[n]
	if !ok {
		return nil, fmt.Errorf("ai: driver [%s] not configured", n)
	}
	return d, nil
}

// Chat runs a chat completion on the default (or named) driver.
func (m *Manager) Chat(req ChatRequest, driver ...string) (*ChatResponse, error) {
	d, err := m.Driver(driver...)
	if err != nil {
		return nil, err
	}
	if req.Model == "" {
		req.Model = "zatrano-fake-1"
	}
	return d.Chat(req)
}

// FakeDriver returns deterministic stub replies.
type FakeDriver struct{}

func (FakeDriver) Name() string { return "fake" }

func (FakeDriver) Chat(req ChatRequest) (*ChatResponse, error) {
	prompt := lastUser(req.Messages)
	reply := "ZATRANO AI stub: " + prompt
	if prompt == "" {
		reply = "ZATRANO AI stub: hello"
	}
	promptTokens := len(strings.Fields(prompt)) + 1
	completionTokens := len(strings.Fields(reply))
	return &ChatResponse{
		ID:      "chat_" + uuid.New()[:8],
		Model:   req.Model,
		Message: Message{Role: "assistant", Content: reply},
		Usage: Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
		Created: time.Now().UTC(),
	}, nil
}

// LogDriver mirrors FakeDriver (placeholder for logging integrations).
type LogDriver struct{}

func (LogDriver) Name() string { return "log" }

func (d LogDriver) Chat(req ChatRequest) (*ChatResponse, error) {
	return FakeDriver{}.Chat(req)
}

func lastUser(messages []Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			return messages[i].Content
		}
	}
	if len(messages) > 0 {
		return messages[len(messages)-1].Content
	}
	return ""
}
