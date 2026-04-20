package notifications

import (
	"context"
	"fmt"

	"github.com/zatrano/zatrano/pkg/mail"
)

// MailChannel sends notifications via the mail manager.
type MailChannel struct {
	manager *mail.Manager
}

// NewMailChannel creates a new email notification channel.
func NewMailChannel(manager *mail.Manager) *MailChannel {
	return &MailChannel{manager: manager}
}

// Name implements Channel.
func (c *MailChannel) Name() string {
	return "mail"
}

// Send implements Channel.
func (c *MailChannel) Send(ctx context.Context, notif Notification) error {
	msg := &mail.Message{
		To:       []mail.Address{{Email: notif.Recipient()}},
		Subject:  notif.Subject(),
		TextBody: notif.Body(),
		HTMLBody: notif.Body(),
	}

	if err := c.manager.Send(ctx, msg); err != nil {
		return fmt.Errorf("mail channel: %w", err)
	}

	return nil
}
