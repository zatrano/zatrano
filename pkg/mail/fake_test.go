package mail

import (
	"testing"
)

func TestFakeMailer_Send(t *testing.T) {
	fake := NewFakeMailer()

	err := fake.Send([]string{"test@example.com"}, "Subject", "Body", false)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	emails := fake.GetSentEmails()
	if len(emails) != 1 {
		t.Errorf("expected 1 email, got %d", len(emails))
	}

	if emails[0].To[0] != "test@example.com" {
		t.Errorf("expected to test@example.com, got %s", emails[0].To[0])
	}
}

func TestFakeMailer_AssertEmailSent(t *testing.T) {
	fake := NewFakeMailer()
	fake.Send([]string{"user@example.com"}, "Welcome", "Hello!", false)

	// Should not panic
	fake.AssertEmailSent("user@example.com", "Welcome")
}

func TestFakeMailer_AssertNoEmailSent(t *testing.T) {
	fake := NewFakeMailer()

	// Should not panic
	fake.AssertNoEmailSent()
}
