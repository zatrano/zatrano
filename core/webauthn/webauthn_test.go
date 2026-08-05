package webauthn_test

import (
	"testing"

	"github.com/zatrano/framework/core/webauthn"
)

func TestRegistrationAndLogin(t *testing.T) {
	m := webauthn.New("localhost", "ZATRANO")
	opts, err := m.BeginRegistration("1", "admin@zatrano.test", "Admin")
	if err != nil {
		t.Fatal(err)
	}
	cred, err := m.FinishRegistration(opts.ChallengeID, "cred-1", "pk")
	if err != nil {
		t.Fatal(err)
	}
	if cred.ID != "cred-1" {
		t.Fatalf("unexpected %#v", cred)
	}
	req, err := m.BeginLogin("1")
	if err != nil {
		t.Fatal(err)
	}
	ok, err := m.FinishLogin(req.ChallengeID, "cred-1")
	if err != nil || !ok {
		t.Fatalf("login failed: %v", err)
	}
	if len(m.CredentialsFor("1")) != 1 {
		t.Fatal("expected one credential")
	}
}
