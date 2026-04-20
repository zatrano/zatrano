package notifications

import (
	"context"
	"fmt"
	"strings"

	"github.com/zatrano/zatrano/pkg/mail"
)

// MailChannel sends notifications via the mail manager.
// TextBody is taken from Notification.Body(). If Data() contains a non-empty
// string value for key "html", it is used as HTMLBody; otherwise HTMLBody
// matches the plain text body. Messages are submitted via Manager.Queue
// (async when a queue is configured on the mail manager).
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
	text := strings.TrimSpace(notif.Body())
	html := text
	if d := notif.Data(); d != nil {
		if v, ok := d["html"].(string); ok {
			if s := strings.TrimSpace(v); s != "" {
				html = s
			}
		}
	}

	msg := &mail.Message{
		To:       []mail.Address{{Email: notif.Recipient()}},
		Subject:  notif.Subject(),
		TextBody: text,
		HTMLBody: html,
	}

	if err := c.manager.Queue(ctx, msg); err != nil {
		return fmt.Errorf("mail channel: %w", err)
	}

	return nil
}
