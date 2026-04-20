package mail

import (
	"fmt"
	"sync"
)

// FakeMailer captures sent emails in memory for testing.
type FakeMailer struct {
	mu     sync.RWMutex
	emails []*SentEmail
}

// SentEmail represents a captured email.
type SentEmail struct {
	To      []string
	Subject string
	Body    string
	HTML    bool
}

// NewFakeMailer creates a new fake mailer.
func NewFakeMailer() *FakeMailer {
	return &FakeMailer{
		emails: make([]*SentEmail, 0),
	}
}

// Send captures the email instead of sending it.
func (f *FakeMailer) Send(to []string, subject, body string, html bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	email := &SentEmail{
		To:      to,
		Subject: subject,
		Body:    body,
		HTML:    html,
	}
	f.emails = append(f.emails, email)
	return nil
}

// GetSentEmails returns all captured emails.
func (f *FakeMailer) GetSentEmails() []*SentEmail {
	f.mu.RLock()
	defer f.mu.RUnlock()

	emails := make([]*SentEmail, len(f.emails))
	copy(emails, f.emails)
	return emails
}

// ClearEmails clears all captured emails.
func (f *FakeMailer) ClearEmails() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.emails = nil
}

// AssertEmailSent asserts that at least one email was sent with the given criteria.
func (f *FakeMailer) AssertEmailSent(to string, subject string) {
	emails := f.GetSentEmails()
	for _, email := range emails {
		if contains(email.To, to) && email.Subject == subject {
			return
		}
	}
	panic("expected email not found")
}

// AssertNoEmailSent asserts that no emails were sent.
func (f *FakeMailer) AssertNoEmailSent() {
	emails := f.GetSentEmails()
	if len(emails) > 0 {
		panic("expected no emails to be sent")
	}
}

// AssertEmailCount asserts the number of emails sent.
func (f *FakeMailer) AssertEmailCount(count int) {
	emails := f.GetSentEmails()
	if len(emails) != count {
		panic(fmt.Sprintf("expected %d emails, got %d", count, len(emails)))
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
