package search

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zatrano/zatrano/pkg/config"
)

// Client wraps a Driver with index naming from configuration.
type Client struct {
	Driver      Driver
	IndexPrefix string
}

// NewClient returns a search client when search.enabled and a supported driver is configured.
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg == nil || !cfg.Search.Enabled {
		return nil, nil
	}
	d := strings.ToLower(strings.TrimSpace(cfg.Search.Driver))
	var drv Driver
	switch d {
	case "meilisearch":
		drv = newMeilisearchDriver(cfg.Search.MeilisearchURL, cfg.Search.MeilisearchAPIKey, &http.Client{Timeout: 60 * time.Second})
	case "typesense":
		drv = newTypesenseDriver(cfg.Search.TypesenseURL, cfg.Search.TypesenseAPIKey, &http.Client{Timeout: 60 * time.Second})
	default:
		return nil, fmt.Errorf("search: unsupported driver %q", cfg.Search.Driver)
	}
	wrapped := &prefixDriver{d: drv, prefix: cfg.Search.DefaultIndexPrefix}
	return &Client{
		Driver:      wrapped,
		IndexPrefix: cfg.Search.DefaultIndexPrefix,
	}, nil
}

// NewDriverForCLI builds a Driver from config (used by zatrano search import without Fiber app).
func NewDriverForCLI(cfg *config.Config) (Driver, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil config")
	}
	c, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}
	if c == nil || c.Driver == nil {
		return nil, fmt.Errorf("search is disabled or not configured; set search.enabled and search.driver for import")
	}
	return c.Driver, nil
}

type prefixDriver struct {
	d      Driver
	prefix string
}

func (p *prefixDriver) IndexName(logical string) string {
	return p.d.IndexName(p.prefix + logical)
}

func (p *prefixDriver) UpsertDocuments(ctx context.Context, logicalIndex string, docs []Document) error {
	return p.d.UpsertDocuments(ctx, p.prefix+logicalIndex, docs)
}

func (p *prefixDriver) DeleteDocument(ctx context.Context, logicalIndex, id string) error {
	return p.d.DeleteDocument(ctx, p.prefix+logicalIndex, id)
}
