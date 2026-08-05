package notification

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zatrano/framework/core/mail"
)

// Notifiable can receive notifications.
type Notifiable interface {
	RouteNotificationFor(channel string) string
	NotificationID() any
}

// Notification is a message that can be sent through channels.
type Notification interface {
	Via() []string
	ToMail(notifiable Notifiable) *mail.Message
	ToDatabase(notifiable Notifiable) map[string]any
	ToBroadcast(notifiable Notifiable) map[string]any
}

// Channel sends a notification through a transport.
type Channel interface {
	Send(notifiable Notifiable, notification Notification) error
}

// Manager sends notifications through registered channels.
type Manager struct {
	channels map[string]Channel
}

// NewManager creates a notification manager.
func NewManager() *Manager {
	return &Manager{channels: make(map[string]Channel)}
}

// Extend registers a channel.
func (m *Manager) Extend(name string, channel Channel) {
	m.channels[name] = channel
}

// Send sends a notification to a notifiable.
func (m *Manager) Send(notifiable Notifiable, notification Notification) error {
	for _, name := range notification.Via() {
		channel, ok := m.channels[name]
		if !ok {
			return fmt.Errorf("notification channel [%s] is not defined", name)
		}
		if err := channel.Send(notifiable, notification); err != nil {
			return err
		}
	}
	return nil
}

// MailChannel sends notifications via mail.
type MailChannel struct {
	mailer *mail.Manager
}

// NewMailChannel creates a mail notification channel.
func NewMailChannel(mailer *mail.Manager) *MailChannel {
	return &MailChannel{mailer: mailer}
}

// Send delivers the mail representation.
func (c *MailChannel) Send(notifiable Notifiable, notification Notification) error {
	message := notification.ToMail(notifiable)
	if message == nil {
		return nil
	}
	if len(message.To) == 0 {
		route := notifiable.RouteNotificationFor("mail")
		if route != "" {
			message.To = []string{route}
		}
	}
	return c.mailer.Send(message)
}

// DatabaseChannel stores notifications in the database.
type DatabaseChannel struct {
	db    *sql.DB
	table string
}

// NewDatabaseChannel creates a database notification channel.
func NewDatabaseChannel(db *sql.DB, table string) *DatabaseChannel {
	if table == "" {
		table = "notifications"
	}
	return &DatabaseChannel{db: db, table: table}
}

// Send stores the notification payload.
func (c *DatabaseChannel) Send(notifiable Notifiable, notification Notification) error {
	payload := notification.ToDatabase(notifiable)
	if payload == nil {
		payload = map[string]any{}
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = c.db.Exec(
		fmt.Sprintf(`INSERT INTO %s (notifiable_id, type, data, created_at) VALUES (?, ?, ?, ?)`, c.table),
		fmt.Sprint(notifiable.NotificationID()),
		fmt.Sprintf("%T", notification),
		string(raw),
		time.Now().Format("2006-01-02 15:04:05"),
	)
	return err
}

// BroadcastChannel publishes notifications to a broadcaster.
type BroadcastChannel struct {
	broadcaster Broadcaster
}

// Broadcaster publishes events to channels.
type Broadcaster interface {
	Broadcast(channel string, event string, payload map[string]any) error
}

// NewBroadcastChannel creates a broadcast notification channel.
func NewBroadcastChannel(broadcaster Broadcaster) *BroadcastChannel {
	return &BroadcastChannel{broadcaster: broadcaster}
}

// Send broadcasts the notification.
func (c *BroadcastChannel) Send(notifiable Notifiable, notification Notification) error {
	payload := notification.ToBroadcast(notifiable)
	if payload == nil {
		payload = map[string]any{}
	}
	channel := notifiable.RouteNotificationFor("broadcast")
	if channel == "" {
		channel = fmt.Sprintf("users.%v", notifiable.NotificationID())
	}
	return c.broadcaster.Broadcast(channel, fmt.Sprintf("%T", notification), payload)
}

// Base can be embedded to provide empty channel representations.
type Base struct{}

// ToMail returns nil by default.
func (Base) ToMail(Notifiable) *mail.Message { return nil }

// ToDatabase returns nil by default.
func (Base) ToDatabase(Notifiable) map[string]any { return nil }

// ToBroadcast returns nil by default.
func (Base) ToBroadcast(Notifiable) map[string]any { return nil }
