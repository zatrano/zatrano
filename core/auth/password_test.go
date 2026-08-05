package auth_test

import (
	"testing"
	"time"

	"github.com/zatrano/framework/core/auth"
)

type stubPasswordProvider struct {
	users map[string]*auth.GenericUser
}

func (p *stubPasswordProvider) RetrieveByID(id any) (auth.Authenticatable, error) {
	return nil, nil
}

func (p *stubPasswordProvider) RetrieveByCredentials(credentials map[string]string) (auth.Authenticatable, error) {
	email := credentials["email"]
	if u, ok := p.users[email]; ok {
		return u, nil
	}
	return nil, nil
}

func (p *stubPasswordProvider) ValidateCredentials(user auth.Authenticatable, credentials map[string]string) bool {
	return true
}

func (p *stubPasswordProvider) UpdatePassword(email, hashedPassword string) error {
	if u, ok := p.users[email]; ok {
		u.Attributes["password"] = hashedPassword
	}
	return nil
}

func TestPasswordBrokerReset(t *testing.T) {
	provider := &stubPasswordProvider{users: map[string]*auth.GenericUser{
		"ada@zatrano.test": {Attributes: map[string]any{"id": 1, "email": "ada@zatrano.test", "password": "old"}},
	}}
	tokens := auth.NewMemoryTokenRepository(time.Hour)
	broker := auth.NewPasswordBroker(tokens, provider, time.Hour)

	token, err := broker.CreateToken("ada@zatrano.test")
	if err != nil {
		t.Fatal(err)
	}
	if !broker.TokenValid("ada@zatrano.test", token) {
		t.Fatal("expected valid token")
	}
	if err := broker.Reset("ada@zatrano.test", token, "secret123"); err != nil {
		t.Fatal(err)
	}
	if broker.TokenValid("ada@zatrano.test", token) {
		t.Fatal("token should be consumed")
	}
}

func TestEmailVerificationHelpers(t *testing.T) {
	user := &auth.GenericUser{Attributes: map[string]any{"email": "ada@zatrano.test"}}
	if auth.HasVerifiedEmail(user) {
		t.Fatal("expected unverified")
	}
	auth.MarkEmailVerified(user.Attributes)
	if !auth.HasVerifiedEmail(user) {
		t.Fatal("expected verified")
	}
	if auth.EmailHash("ada@zatrano.test") == "" {
		t.Fatal("expected hash")
	}
}
