package notification

import (
	"fmt"
	"sync"
)

// PushNotification optionally supplies a push payload.
type PushNotification interface {
	ToPush(notifiable Notifiable) map[string]any
}

// PushSender delivers push payloads to a device token.
type PushSender interface {
	Send(deviceToken string, payload map[string]any) error
}

// MemoryPushSender records push deliveries for tests/demos.
type MemoryPushSender struct {
	mu      sync.Mutex
	Entries []PushEntry
}

// PushEntry is a recorded push delivery.
type PushEntry struct {
	Token   string
	Payload map[string]any
}

// Send records the push.
func (s *MemoryPushSender) Send(deviceToken string, payload map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Entries = append(s.Entries, PushEntry{Token: deviceToken, Payload: payload})
	return nil
}

// Last returns the most recent entry.
func (s *MemoryPushSender) Last() (PushEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.Entries) == 0 {
		return PushEntry{}, false
	}
	return s.Entries[len(s.Entries)-1], true
}

// PushChannel sends notifications via a push sender stub.
type PushChannel struct {
	sender PushSender
}

// NewPushChannel creates a push notification channel.
func NewPushChannel(sender PushSender) *PushChannel {
	if sender == nil {
		sender = &MemoryPushSender{}
	}
	return &PushChannel{sender: sender}
}

// Send delivers the push representation.
func (c *PushChannel) Send(notifiable Notifiable, notification Notification) error {
	var payload map[string]any
	if p, ok := notification.(PushNotification); ok {
		payload = p.ToPush(notifiable)
	}
	if payload == nil {
		payload = notification.ToBroadcast(notifiable)
	}
	if payload == nil {
		payload = map[string]any{"type": fmt.Sprintf("%T", notification)}
	}
	token := notifiable.RouteNotificationFor("push")
	if token == "" {
		token = fmt.Sprintf("user:%v", notifiable.NotificationID())
	}
	return c.sender.Send(token, payload)
}
