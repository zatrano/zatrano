package notifications

import (
	"context"
	"time"
)

// Notification represents a notification to be sent through one or more channels.
type Notification interface {
	// Subject returns the notification subject (e.g., email subject).
	Subject() string
	// Body returns the notification body/content.
	Body() string
	// Recipient returns the target recipient (email, phone, user ID, etc.).
	Recipient() string
	// Data returns optional structured data for the notification.
	Data() map[string]any
}

// Channel is a notification delivery mechanism.
type Channel interface {
	// Name returns the unique name of this channel (e.g., "mail", "sms", "database", "push").
	Name() string
	// Send delivers the notification through this channel.
	Send(ctx context.Context, notif Notification) error
}

// Manager is responsible for routing notifications to multiple channels.
type Manager struct {
	channels map[string]Channel
}

// NewManager creates a new notification manager.
func NewManager() *Manager {
	return &Manager{channels: make(map[string]Channel)}
}

// Register adds a channel to the manager.
func (m *Manager) Register(channel Channel) {
	m.channels[channel.Name()] = channel
}

// Send dispatches the notification to all registered channels.
func (m *Manager) Send(ctx context.Context, notif Notification) error {
	for _, ch := range m.channels {
		if err := ch.Send(ctx, notif); err != nil {
			// Log error but continue sending to other channels
			_ = err
		}
	}
	return nil
}

// SendToChannels dispatches the notification to specific channels.
func (m *Manager) SendToChannels(ctx context.Context, notif Notification, channelNames ...string) error {
	for _, name := range channelNames {
		if ch, ok := m.channels[name]; ok {
			if err := ch.Send(ctx, notif); err != nil {
				_ = err
			}
		}
	}
	return nil
}

// BaseNotification is a simple notification implementation.
type BaseNotification struct {
	subject   string
	body      string
	recipient string
	data      map[string]any
}

// NewNotification creates a new base notification.
func NewNotification(subject, body, recipient string) *BaseNotification {
	return &BaseNotification{
		subject:   subject,
		body:      body,
		recipient: recipient,
		data:      make(map[string]any),
	}
}

// Subject implements Notification.
func (n *BaseNotification) Subject() string {
	return n.subject
}

// Body implements Notification.
func (n *BaseNotification) Body() string {
	return n.body
}

// Recipient implements Notification.
func (n *BaseNotification) Recipient() string {
	return n.recipient
}

// Data implements Notification.
func (n *BaseNotification) Data() map[string]any {
	return n.data
}

// WithData adds structured data to the notification.
func (n *BaseNotification) WithData(key string, value any) *BaseNotification {
	n.data[key] = value
	return n
}

// NotificationRecord represents a stored notification in the database.
type NotificationRecord struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `json:"user_id"`
	Type      string     `json:"type"`
	Subject   string     `json:"subject"`
	Body      string     `json:"body"`
	Data      string     `json:"data"` // JSON serialized
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName specifies the table name for NotificationRecord.
func (NotificationRecord) TableName() string {
	return "notifications"
}
