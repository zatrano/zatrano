package notifications

import (
	"context"
	"testing"

	"github.com/zatrano/zatrano/pkg/mail"
	"go.uber.org/zap"
)

func TestNotificationManager_Send(t *testing.T) {
	manager := NewManager()
	callCount := 0

	// Register a mock channel
	manager.Register(mockChannel{
		name: "test",
		send: func(ctx context.Context, notif Notification) error {
			callCount++
			return nil
		},
	})

	notif := NewNotification("Subject", "Body", "recipient@example.com")
	ctx := context.Background()
	if err := manager.Send(ctx, notif); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if callCount != 1 {
		t.Fatalf("expected 1 call, got %d", callCount)
	}
}

func TestNotificationManager_SendToChannels(t *testing.T) {
	manager := NewManager()
	var called []string

	manager.Register(mockChannel{
		name: "channel1",
		send: func(ctx context.Context, notif Notification) error {
			called = append(called, "channel1")
			return nil
		},
	})
	manager.Register(mockChannel{
		name: "channel2",
		send: func(ctx context.Context, notif Notification) error {
			called = append(called, "channel2")
			return nil
		},
	})

	notif := NewNotification("Subject", "Body", "recipient@example.com")
	ctx := context.Background()
	if err := manager.SendToChannels(ctx, notif, "channel1"); err != nil {
		t.Fatalf("SendToChannels failed: %v", err)
	}

	if len(called) != 1 || called[0] != "channel1" {
		t.Fatalf("expected [channel1], got %v", called)
	}
}

func TestMailChannel_Send(t *testing.T) {
	sent := &mail.Message{}
	driver := &mockMailDriver{sent: sent}
	manager := mail.New(driver, mail.MailConfig{FromName: "Zatrano", FromEmail: "noreply@example.com"}, zap.NewNop(), nil)
	channel := NewMailChannel(manager)

	notif := NewNotification("Welcome", "Welcome to the platform!", "user@example.com")
	if err := channel.Send(context.Background(), notif); err != nil {
		t.Fatalf("MailChannel.Send failed: %v", err)
	}

	if len(sent.To) != 1 || sent.To[0].Email != "user@example.com" {
		t.Fatalf("expected recipient user@example.com, got %v", sent.To)
	}
	if sent.Subject != "Welcome" {
		t.Fatalf("expected subject Welcome, got %q", sent.Subject)
	}
	if sent.TextBody != "Welcome to the platform!" {
		t.Fatalf("expected text body to match notification body")
	}
	if sent.HTMLBody != "Welcome to the platform!" {
		t.Fatalf("expected html body to match notification body")
	}
}

func TestMailChannel_Send_HTMLData(t *testing.T) {
	sent := &mail.Message{}
	driver := &mockMailDriver{sent: sent}
	manager := mail.New(driver, mail.MailConfig{FromName: "Z", FromEmail: "z@example.com"}, zap.NewNop(), nil)
	channel := NewMailChannel(manager)

	notif := NewNotification("Sub", "plain text", "u@example.com").WithData("html", "<p>rich</p>")
	if err := channel.Send(context.Background(), notif); err != nil {
		t.Fatalf("MailChannel.Send failed: %v", err)
	}
	if sent.TextBody != "plain text" {
		t.Fatalf("expected text body plain text, got %q", sent.TextBody)
	}
	if sent.HTMLBody != "<p>rich</p>" {
		t.Fatalf("expected html from data, got %q", sent.HTMLBody)
	}
}

func TestBaseNotification_WithData(t *testing.T) {
	notif := NewNotification("Subject", "Body", "recipient").
		WithData("key1", "value1").
		WithData("key2", 123)

	data := notif.Data()
	if data["key1"] != "value1" || data["key2"] != 123 {
		t.Fatalf("expected data to contain key1 and key2")
	}
}

type mockChannel struct {
	name string
	send func(context.Context, Notification) error
}

func (m mockChannel) Name() string {
	return m.name
}

func (m mockChannel) Send(ctx context.Context, notif Notification) error {
	return m.send(ctx, notif)
}

type mockMailDriver struct {
	sent *mail.Message
	err  error
}

func (m *mockMailDriver) Send(ctx context.Context, msg *mail.Message) error {
	if m.err != nil {
		return m.err
	}
	*m.sent = *msg
	return nil
}

func (m *mockMailDriver) Name() string {
	return "mock"
}
