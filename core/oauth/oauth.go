package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/zatrano/framework/core/http"
	"github.com/zatrano/framework/core/routing"
	"github.com/zatrano/framework/core/support/uuid"
)

// Client is an OAuth2 client application.
type Client struct {
	ID           string   `json:"id"`
	Secret       string   `json:"secret,omitempty"`
	Name         string   `json:"name"`
	RedirectURIs []string `json:"redirect_uris"`
	Scopes       []string `json:"scopes"`
}

// AccessToken is an issued OAuth2 token.
type AccessToken struct {
	Token     string    `json:"access_token"`
	Type      string    `json:"token_type"`
	ExpiresIn int       `json:"expires_in"`
	Scope     string    `json:"scope,omitempty"`
	ClientID  string    `json:"client_id"`
	UserID    string    `json:"user_id,omitempty"`
	ExpiresAt time.Time `json:"-"`
}

// Server is an in-memory OAuth2 authorization server stub.
type Server struct {
	mu      sync.Mutex
	clients map[string]*Client
	codes   map[string]authCode
	tokens  map[string]*AccessToken
	ttl     time.Duration
}

type authCode struct {
	ClientID    string
	RedirectURI string
	UserID      string
	Scope       string
	ExpiresAt   time.Time
}

// New creates an OAuth2 server.
func New() *Server {
	return &Server{
		clients: make(map[string]*Client),
		codes:   make(map[string]authCode),
		tokens:  make(map[string]*AccessToken),
		ttl:     time.Hour,
	}
}

// RegisterClient stores a client (generates ID/secret when empty).
func (s *Server) RegisterClient(c Client) *Client {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c.ID == "" {
		c.ID = "client_" + uuid.New()[:8]
	}
	if c.Secret == "" {
		c.Secret = randomToken(24)
	}
	if len(c.Scopes) == 0 {
		c.Scopes = []string{"*"}
	}
	cp := c
	s.clients[c.ID] = &cp
	return &cp
}

// Client returns a registered client.
func (s *Server) Client(id string) (*Client, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.clients[id]
	if !ok {
		return nil, fmt.Errorf("oauth: unknown client")
	}
	cp := *c
	cp.Secret = ""
	return &cp, nil
}

// Authorize creates an authorization code (authorization_code grant).
func (s *Server) Authorize(clientID, redirectURI, userID, scope string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	client, ok := s.clients[clientID]
	if !ok {
		return "", fmt.Errorf("oauth: unknown client")
	}
	if !validRedirect(client.RedirectURIs, redirectURI) {
		return "", fmt.Errorf("oauth: invalid redirect_uri")
	}
	if scope == "" {
		scope = strings.Join(client.Scopes, " ")
	}
	code := randomToken(32)
	s.codes[code] = authCode{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		UserID:      userID,
		Scope:       scope,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}
	return code, nil
}

// Token exchanges a grant for an access token.
func (s *Server) Token(grantType, clientID, clientSecret string, params map[string]string) (*AccessToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	client, ok := s.clients[clientID]
	if !ok || client.Secret != clientSecret {
		return nil, fmt.Errorf("oauth: invalid client credentials")
	}

	switch grantType {
	case "client_credentials":
		scope := params["scope"]
		if scope == "" {
			scope = strings.Join(client.Scopes, " ")
		}
		return s.issueLocked(clientID, "", scope), nil
	case "authorization_code":
		code := params["code"]
		redirect := params["redirect_uri"]
		entry, ok := s.codes[code]
		if !ok || time.Now().After(entry.ExpiresAt) {
			return nil, fmt.Errorf("oauth: invalid authorization code")
		}
		if entry.ClientID != clientID || entry.RedirectURI != redirect {
			return nil, fmt.Errorf("oauth: code mismatch")
		}
		delete(s.codes, code)
		return s.issueLocked(clientID, entry.UserID, entry.Scope), nil
	default:
		return nil, fmt.Errorf("oauth: unsupported grant_type")
	}
}

// Introspect validates an access token.
func (s *Server) Introspect(token string) map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tokens[hashToken(token)]
	if !ok || time.Now().After(t.ExpiresAt) {
		return map[string]any{"active": false}
	}
	return map[string]any{
		"active":     true,
		"client_id":  t.ClientID,
		"user_id":    t.UserID,
		"scope":      t.Scope,
		"exp":        t.ExpiresAt.Unix(),
		"token_type": "Bearer",
	}
}

// AuthorizeHandler handles GET /oauth/authorize.
func (s *Server) AuthorizeHandler() routing.HandlerFunc {
	return func(req *http.Request) *http.Response {
		clientID := req.Query("client_id")
		redirect := req.Query("redirect_uri")
		userID := req.Query("user_id", "1")
		scope := req.Query("scope")
		state := req.Query("state")
		code, err := s.Authorize(clientID, redirect, userID, scope)
		if err != nil {
			return http.JSON(map[string]any{"error": "invalid_request", "error_description": err.Error()}).Status(400)
		}
		if redirect == "" {
			return http.JSON(map[string]any{"code": code, "state": state})
		}
		u, err := url.Parse(redirect)
		if err != nil {
			return http.JSON(map[string]any{"code": code, "state": state})
		}
		q := u.Query()
		q.Set("code", code)
		if state != "" {
			q.Set("state", state)
		}
		u.RawQuery = q.Encode()
		return http.Redirect(u.String(), 302)
	}
}

// TokenHandler handles POST /oauth/token.
func (s *Server) TokenHandler() routing.HandlerFunc {
	return func(req *http.Request) *http.Response {
		grant := req.Input("grant_type")
		clientID := req.Input("client_id")
		secret := req.Input("client_secret")
		if clientID == "" {
			clientID, secret, _ = parseBasicAuth(req.Header("Authorization"))
		}
		token, err := s.Token(grant, clientID, secret, map[string]string{
			"code":         req.Input("code"),
			"redirect_uri": req.Input("redirect_uri"),
			"scope":        req.Input("scope"),
		})
		if err != nil {
			return http.JSON(map[string]any{"error": "invalid_grant", "error_description": err.Error()}).Status(400)
		}
		return http.JSON(token)
	}
}

// IntrospectHandler handles POST /oauth/introspect.
func (s *Server) IntrospectHandler() routing.HandlerFunc {
	return func(req *http.Request) *http.Response {
		return http.JSON(s.Introspect(req.Input("token")))
	}
}

func (s *Server) issueLocked(clientID, userID, scope string) *AccessToken {
	plain := randomToken(40)
	token := &AccessToken{
		Token:     plain,
		Type:      "Bearer",
		ExpiresIn: int(s.ttl.Seconds()),
		Scope:     scope,
		ClientID:  clientID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.ttl),
	}
	s.tokens[hashToken(plain)] = token
	return token
}

func validRedirect(allowed []string, redirect string) bool {
	if redirect == "" {
		return true
	}
	for _, a := range allowed {
		if a == redirect || a == "*" {
			return true
		}
	}
	return false
}

func randomToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func parseBasicAuth(header string) (string, string, bool) {
	if !strings.HasPrefix(header, "Basic ") {
		return "", "", false
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(header, "Basic "))
	if err != nil {
		return "", "", false
	}
	parts := strings.SplitN(string(raw), ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}
