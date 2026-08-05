package mongo

import (
	"fmt"
	"strings"
	"sync"

	"github.com/zatrano/framework/core/support/uuid"
)

// Client is an in-memory MongoDB-like stub.
type Client struct {
	mu        sync.Mutex
	URI       string
	databases map[string]*Database
}

// Database holds named collections.
type Database struct {
	mu          sync.Mutex
	Name        string
	collections map[string]*Collection
}

// Collection stores documents.
type Collection struct {
	mu   sync.Mutex
	Name string
	docs []map[string]any
}

// Connect opens a stub client (any URI works; "memory" is the default driver).
func Connect(uri string) *Client {
	if uri == "" {
		uri = "memory"
	}
	return &Client{
		URI:       uri,
		databases: make(map[string]*Database),
	}
}

// Database returns or creates a database.
func (c *Client) Database(name string) *Database {
	c.mu.Lock()
	defer c.mu.Unlock()
	if name == "" {
		name = "zatrano"
	}
	db, ok := c.databases[name]
	if !ok {
		db = &Database{Name: name, collections: make(map[string]*Collection)}
		c.databases[name] = db
	}
	return db
}

// Ping reports the stub is reachable.
func (c *Client) Ping() error {
	if c == nil {
		return fmt.Errorf("mongo: nil client")
	}
	return nil
}

// Collection returns or creates a collection.
func (d *Database) Collection(name string) *Collection {
	d.mu.Lock()
	defer d.mu.Unlock()
	if name == "" {
		name = "items"
	}
	col, ok := d.collections[name]
	if !ok {
		col = &Collection{Name: name, docs: make([]map[string]any, 0)}
		d.collections[name] = col
	}
	return col
}

// InsertOne inserts a document and returns its _id.
func (c *Collection) InsertOne(doc map[string]any) (string, error) {
	if doc == nil {
		return "", fmt.Errorf("mongo: document is required")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	cloned := cloneDoc(doc)
	id, _ := cloned["_id"].(string)
	if id == "" {
		id = uuid.New()
		cloned["_id"] = id
	}
	c.docs = append(c.docs, cloned)
	return id, nil
}

// Find returns documents matching an equality filter (empty = all).
func (c *Collection) Find(filter map[string]any) ([]map[string]any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]map[string]any, 0)
	for _, doc := range c.docs {
		if match(doc, filter) {
			out = append(out, cloneDoc(doc))
		}
	}
	return out, nil
}

// FindOne returns the first matching document.
func (c *Collection) FindOne(filter map[string]any) (map[string]any, error) {
	docs, err := c.Find(filter)
	if err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, fmt.Errorf("mongo: no documents")
	}
	return docs[0], nil
}

// UpdateOne sets fields on the first matching document.
func (c *Collection) UpdateOne(filter, update map[string]any) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, doc := range c.docs {
		if match(doc, filter) {
			for k, v := range update {
				if k == "_id" {
					continue
				}
				doc[k] = v
			}
			c.docs[i] = doc
			return true, nil
		}
	}
	return false, nil
}

// DeleteOne removes the first matching document.
func (c *Collection) DeleteOne(filter map[string]any) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, doc := range c.docs {
		if match(doc, filter) {
			c.docs = append(c.docs[:i], c.docs[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

// Count returns document count.
func (c *Collection) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.docs)
}

func match(doc, filter map[string]any) bool {
	if len(filter) == 0 {
		return true
	}
	for k, want := range filter {
		got, ok := doc[k]
		if !ok {
			return false
		}
		if stringify(got) != stringify(want) {
			return false
		}
	}
	return true
}

func stringify(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	default:
		return fmt.Sprint(v)
	}
}

func cloneDoc(doc map[string]any) map[string]any {
	out := make(map[string]any, len(doc))
	for k, v := range doc {
		out[k] = v
	}
	return out
}

// ParseDatabase extracts a database name hint from a mongodb URI.
func ParseDatabase(uri string) string {
	uri = strings.TrimSpace(uri)
	if uri == "" || uri == "memory" {
		return "zatrano"
	}
	if idx := strings.LastIndex(uri, "/"); idx >= 0 && idx+1 < len(uri) {
		rest := uri[idx+1:]
		if q := strings.IndexAny(rest, "?#"); q >= 0 {
			rest = rest[:q]
		}
		if rest != "" {
			return rest
		}
	}
	return "zatrano"
}
