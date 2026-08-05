package billing

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/zatrano/framework/core/support/uuid"
)

// Customer is a billable account.
type Customer struct {
	ID                   string    `json:"id"`
	Email                string    `json:"email"`
	Name                 string    `json:"name,omitempty"`
	DefaultPaymentMethod string    `json:"default_payment_method,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

// Subscription is a recurring plan attachment.
type Subscription struct {
	ID         string     `json:"id"`
	CustomerID string     `json:"customer_id"`
	Name       string     `json:"name"`
	PriceID    string     `json:"price_id"`
	Status     string     `json:"status"`
	TrialEnds  *time.Time `json:"trial_ends_at,omitempty"`
	EndsAt     *time.Time `json:"ends_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// Invoice is a one-off charge record.
type Invoice struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Amount     int64     `json:"amount"`
	Currency   string    `json:"currency"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// CheckoutSession is a hosted checkout stub.
type CheckoutSession struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	PriceID    string `json:"price_id"`
	URL        string `json:"url"`
	Status     string `json:"status"`
}

// Manager is an in-memory billing stub (Cashier-like surface).
type Manager struct {
	mu            sync.RWMutex
	baseURL       string
	customers     map[string]*Customer
	subscriptions map[string]*Subscription
	invoices      map[string]*Invoice
	checkouts     map[string]*CheckoutSession
}

// New creates a billing manager.
func New(baseURL string) *Manager {
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &Manager{
		baseURL:       strings.TrimRight(baseURL, "/"),
		customers:     make(map[string]*Customer),
		subscriptions: make(map[string]*Subscription),
		invoices:      make(map[string]*Invoice),
		checkouts:     make(map[string]*CheckoutSession),
	}
}

// CreateCustomer registers a customer.
func (m *Manager) CreateCustomer(email, name string) (*Customer, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, fmt.Errorf("billing: email is required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, c := range m.customers {
		if c.Email == email {
			return c, nil
		}
	}
	c := &Customer{
		ID:        "cus_" + uuid.New()[:8],
		Email:     email,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
	m.customers[c.ID] = c
	return c, nil
}

// Customer returns a customer by ID.
func (m *Manager) Customer(id string) (*Customer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.customers[id]
	if !ok {
		return nil, fmt.Errorf("billing: customer [%s] not found", id)
	}
	return c, nil
}

// Subscribe starts a subscription (optionally with a trial).
func (m *Manager) Subscribe(customerID, name, priceID string, trialDays ...int) (*Subscription, error) {
	if _, err := m.Customer(customerID); err != nil {
		return nil, err
	}
	if name == "" {
		name = "default"
	}
	if priceID == "" {
		return nil, fmt.Errorf("billing: price_id is required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	sub := &Subscription{
		ID:         "sub_" + uuid.New()[:8],
		CustomerID: customerID,
		Name:       name,
		PriceID:    priceID,
		Status:     "active",
		CreatedAt:  time.Now().UTC(),
	}
	if len(trialDays) > 0 && trialDays[0] > 0 {
		t := time.Now().UTC().Add(time.Duration(trialDays[0]) * 24 * time.Hour)
		sub.TrialEnds = &t
		sub.Status = "trialing"
	}
	m.subscriptions[sub.ID] = sub
	return sub, nil
}

// Cancel ends a subscription immediately or at period end.
func (m *Manager) Cancel(subscriptionID string, immediately bool) (*Subscription, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sub, ok := m.subscriptions[subscriptionID]
	if !ok {
		return nil, fmt.Errorf("billing: subscription [%s] not found", subscriptionID)
	}
	now := time.Now().UTC()
	if immediately {
		sub.Status = "canceled"
		sub.EndsAt = &now
	} else {
		sub.Status = "canceling"
		end := now.Add(30 * 24 * time.Hour)
		sub.EndsAt = &end
	}
	return sub, nil
}

// Subscribed reports whether the customer has an active/trialing subscription.
func (m *Manager) Subscribed(customerID, name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if name == "" {
		name = "default"
	}
	for _, sub := range m.subscriptions {
		if sub.CustomerID == customerID && sub.Name == name {
			if sub.Status == "active" || sub.Status == "trialing" || sub.Status == "canceling" {
				if sub.EndsAt != nil && time.Now().After(*sub.EndsAt) {
					continue
				}
				return true
			}
		}
	}
	return false
}

// OnTrial reports whether the customer is currently on trial.
func (m *Manager) OnTrial(customerID, name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if name == "" {
		name = "default"
	}
	now := time.Now().UTC()
	for _, sub := range m.subscriptions {
		if sub.CustomerID == customerID && sub.Name == name && sub.TrialEnds != nil {
			if now.Before(*sub.TrialEnds) && (sub.Status == "trialing" || sub.Status == "active") {
				return true
			}
		}
	}
	return false
}

// Checkout creates a hosted checkout session URL.
func (m *Manager) Checkout(customerID, priceID string) (*CheckoutSession, error) {
	if _, err := m.Customer(customerID); err != nil {
		return nil, err
	}
	if priceID == "" {
		return nil, fmt.Errorf("billing: price_id is required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	session := &CheckoutSession{
		ID:         "cs_" + uuid.New()[:8],
		CustomerID: customerID,
		PriceID:    priceID,
		Status:     "open",
	}
	session.URL = fmt.Sprintf("%s/billing/checkout/%s", m.baseURL, session.ID)
	m.checkouts[session.ID] = session
	return session, nil
}

// Invoice charges a customer once.
func (m *Manager) Invoice(customerID string, amount int64, currency string) (*Invoice, error) {
	if _, err := m.Customer(customerID); err != nil {
		return nil, err
	}
	if amount <= 0 {
		return nil, fmt.Errorf("billing: amount must be positive")
	}
	if currency == "" {
		currency = "usd"
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	inv := &Invoice{
		ID:         "in_" + uuid.New()[:8],
		CustomerID: customerID,
		Amount:     amount,
		Currency:   strings.ToLower(currency),
		Status:     "paid",
		CreatedAt:  time.Now().UTC(),
	}
	m.invoices[inv.ID] = inv
	return inv, nil
}

// SubscriptionsFor returns subscriptions for a customer.
func (m *Manager) SubscriptionsFor(customerID string) []*Subscription {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Subscription, 0)
	for _, sub := range m.subscriptions {
		if sub.CustomerID == customerID {
			out = append(out, sub)
		}
	}
	return out
}
