package notifications

import (
	"context"
	"fmt"
)

// SMSChannel sends notifications via SMS.
type SMSChannel struct {
	driver SMSDriver
}

// SMSDriver is an interface for SMS providers.
type SMSDriver interface {
	Send(ctx context.Context, phoneNumber, message string) error
}

// NewSMSChannel creates an SMS notification channel.
func NewSMSChannel(driver SMSDriver) *SMSChannel {
	return &SMSChannel{driver: driver}
}

// Name implements Channel.
func (c *SMSChannel) Name() string {
	return "sms"
}

// Send implements Channel.
func (c *SMSChannel) Send(ctx context.Context, notif Notification) error {
	return c.driver.Send(ctx, notif.Recipient(), notif.Body())
}

// TwilioDriver is a Twilio SMS driver.
type TwilioDriver struct {
	accountSID string
	authToken  string
	fromNumber string
}

// NewTwilioDriver creates a Twilio SMS driver.
func NewTwilioDriver(accountSID, authToken, fromNumber string) *TwilioDriver {
	return &TwilioDriver{
		accountSID: accountSID,
		authToken:  authToken,
		fromNumber: fromNumber,
	}
}

// Send sends an SMS via Twilio.
func (d *TwilioDriver) Send(ctx context.Context, phoneNumber, message string) error {
	// TODO: Implement Twilio API call
	// This would require github.com/twilio/twilio-go
	_ = ctx
	_ = phoneNumber
	_ = message
	return fmt.Errorf("twilio driver not yet implemented")
}

// NetgsmDriver is a Netgsm SMS driver (Turkish SMS provider).
type NetgsmDriver struct {
	username string
	password string
	senderID string
	apiURL   string
}

// NewNetgsmDriver creates a Netgsm SMS driver.
func NewNetgsmDriver(username, password, senderID string) *NetgsmDriver {
	return &NetgsmDriver{
		username: username,
		password: password,
		senderID: senderID,
		apiURL:   "https://api.netgsm.com.tr",
	}
}

// Send sends an SMS via Netgsm.
func (d *NetgsmDriver) Send(ctx context.Context, phoneNumber, message string) error {
	// TODO: Implement Netgsm API call
	// This would require HTTP POST to Netgsm endpoint
	_ = ctx
	_ = phoneNumber
	_ = message
	return fmt.Errorf("netgsm driver not yet implemented")
}
