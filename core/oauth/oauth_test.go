package oauth_test

import (
	"testing"

	"github.com/zatrano/framework/core/oauth"
)

func TestOAuthClientCredentials(t *testing.T) {
	s := oauth.New()
	client := s.RegisterClient(oauth.Client{
		Name:         "Demo",
		RedirectURIs: []string{"http://localhost/callback"},
		Scopes:       []string{"read", "write"},
	})
	token, err := s.Token("client_credentials", client.ID, client.Secret, nil)
	if err != nil || token.Token == "" {
		t.Fatalf("token=%v err=%v", token, err)
	}
	info := s.Introspect(token.Token)
	if info["active"] != true {
		t.Fatalf("introspect=%v", info)
	}
}

func TestOAuthAuthorizationCode(t *testing.T) {
	s := oauth.New()
	client := s.RegisterClient(oauth.Client{
		Name:         "Web",
		RedirectURIs: []string{"http://localhost/callback"},
	})
	code, err := s.Authorize(client.ID, "http://localhost/callback", "42", "read")
	if err != nil {
		t.Fatal(err)
	}
	token, err := s.Token("authorization_code", client.ID, client.Secret, map[string]string{
		"code":         code,
		"redirect_uri": "http://localhost/callback",
	})
	if err != nil || token.UserID != "42" {
		t.Fatalf("token=%v err=%v", token, err)
	}
}
