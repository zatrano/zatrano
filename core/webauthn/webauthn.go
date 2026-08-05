package webauthn

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Challenge is a pending WebAuthn ceremony.
type Challenge struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"` // registration | authentication
	Challenge string    `json:"challenge"`
	RPID      string    `json:"rp_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CreationOptions is a stub PublicKeyCredentialCreationOptions payload.
type CreationOptions struct {
	Challenge        string           `json:"challenge"`
	RPID             string           `json:"rp_id"`
	RPName           string           `json:"rp_name"`
	UserID           string           `json:"user_id"`
	UserName         string           `json:"user_name"`
	UserDisplayName  string           `json:"user_display_name"`
	TimeoutMS        int              `json:"timeout_ms"`
	ChallengeID      string           `json:"challenge_id"`
	PubKeyCredParams []map[string]any `json:"pub_key_cred_params"`
}

// RequestOptions is a stub PublicKeyCredentialRequestOptions payload.
type RequestOptions struct {
	Challenge   string `json:"challenge"`
	RPID        string `json:"rp_id"`
	TimeoutMS   int    `json:"timeout_ms"`
	ChallengeID string `json:"challenge_id"`
	UserID      string `json:"user_id"`
}

// Credential is a stub registered authenticator.
type Credential struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	PublicKey string    `json:"public_key"`
	CreatedAt time.Time `json:"created_at"`
}

// Manager is an in-memory WebAuthn stub (not production crypto).
type Manager struct {
	mu          sync.RWMutex
	rpID        string
	rpName      string
	ttl         time.Duration
	challenges  map[string]Challenge
	credentials map[string][]Credential // userID -> creds
}

// New creates a WebAuthn stub manager.
func New(rpID, rpName string) *Manager {
	if strings.TrimSpace(rpID) == "" {
		rpID = "localhost"
	}
	if strings.TrimSpace(rpName) == "" {
		rpName = "ZATRANO"
	}
	return &Manager{
		rpID:        rpID,
		rpName:      rpName,
		ttl:         5 * time.Minute,
		challenges:  make(map[string]Challenge),
		credentials: make(map[string][]Credential),
	}
}

// BeginRegistration starts a registration ceremony.
func (m *Manager) BeginRegistration(userID, userName, displayName string) (*CreationOptions, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("webauthn: user id required")
	}
	if userName == "" {
		userName = userID
	}
	if displayName == "" {
		displayName = userName
	}
	ch, err := m.createChallenge(userID, "registration")
	if err != nil {
		return nil, err
	}
	return &CreationOptions{
		Challenge:       ch.Challenge,
		RPID:            m.rpID,
		RPName:          m.rpName,
		UserID:          userID,
		UserName:        userName,
		UserDisplayName: displayName,
		TimeoutMS:       int(m.ttl / time.Millisecond),
		ChallengeID:     ch.ID,
		PubKeyCredParams: []map[string]any{
			{"type": "public-key", "alg": -7},
			{"type": "public-key", "alg": -257},
		},
	}, nil
}

// FinishRegistration stores a stub credential (accepts any non-empty credential id).
func (m *Manager) FinishRegistration(challengeID, credentialID, publicKey string) (*Credential, error) {
	ch, err := m.takeChallenge(challengeID, "registration")
	if err != nil {
		return nil, err
	}
	credentialID = strings.TrimSpace(credentialID)
	if credentialID == "" {
		return nil, fmt.Errorf("webauthn: credential id required")
	}
	if publicKey == "" {
		publicKey = "stub-public-key"
	}
	cred := Credential{
		ID:        credentialID,
		UserID:    ch.UserID,
		PublicKey: publicKey,
		CreatedAt: time.Now().UTC(),
	}
	m.mu.Lock()
	m.credentials[ch.UserID] = append(m.credentials[ch.UserID], cred)
	m.mu.Unlock()
	return &cred, nil
}

// BeginLogin starts an authentication ceremony.
func (m *Manager) BeginLogin(userID string) (*RequestOptions, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("webauthn: user id required")
	}
	m.mu.RLock()
	creds := m.credentials[userID]
	m.mu.RUnlock()
	if len(creds) == 0 {
		return nil, fmt.Errorf("webauthn: no credentials for user")
	}
	ch, err := m.createChallenge(userID, "authentication")
	if err != nil {
		return nil, err
	}
	return &RequestOptions{
		Challenge:   ch.Challenge,
		RPID:        m.rpID,
		TimeoutMS:   int(m.ttl / time.Millisecond),
		ChallengeID: ch.ID,
		UserID:      userID,
	}, nil
}

// FinishLogin verifies a stub assertion (credential must exist for the challenge user).
func (m *Manager) FinishLogin(challengeID, credentialID string) (bool, error) {
	ch, err := m.takeChallenge(challengeID, "authentication")
	if err != nil {
		return false, err
	}
	credentialID = strings.TrimSpace(credentialID)
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, cred := range m.credentials[ch.UserID] {
		if cred.ID == credentialID {
			return true, nil
		}
	}
	return false, fmt.Errorf("webauthn: unknown credential")
}

// CredentialsFor returns registered stub credentials for a user.
func (m *Manager) CredentialsFor(userID string) []Credential {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Credential, len(m.credentials[userID]))
	copy(out, m.credentials[userID])
	return out
}

func (m *Manager) createChallenge(userID, typ string) (Challenge, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return Challenge{}, err
	}
	idRaw := make([]byte, 16)
	if _, err := rand.Read(idRaw); err != nil {
		return Challenge{}, err
	}
	ch := Challenge{
		ID:        base64.RawURLEncoding.EncodeToString(idRaw),
		UserID:    userID,
		Type:      typ,
		Challenge: base64.RawURLEncoding.EncodeToString(raw),
		RPID:      m.rpID,
		ExpiresAt: time.Now().Add(m.ttl),
	}
	m.mu.Lock()
	m.challenges[ch.ID] = ch
	m.mu.Unlock()
	return ch, nil
}

func (m *Manager) takeChallenge(id, typ string) (Challenge, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ch, ok := m.challenges[id]
	if !ok {
		return Challenge{}, fmt.Errorf("webauthn: challenge not found")
	}
	delete(m.challenges, id)
	if time.Now().After(ch.ExpiresAt) {
		return Challenge{}, fmt.Errorf("webauthn: challenge expired")
	}
	if ch.Type != typ {
		return Challenge{}, fmt.Errorf("webauthn: challenge type mismatch")
	}
	return ch, nil
}
