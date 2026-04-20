package search

import (
	"context"
)

// Document is a single row sent to an external search index.
type Document struct {
	ID     string
	Fields map[string]any
}

// Driver indexes documents in Meilisearch or Typesense.
type Driver interface {
	// IndexName returns the physical index/collection name for a logical name.
	IndexName(logical string) string
	// UpsertDocuments sends documents (create or replace by primary key).
	UpsertDocuments(ctx context.Context, logicalIndex string, docs []Document) error
	// DeleteDocument removes one document by id.
	DeleteDocument(ctx context.Context, logicalIndex, id string) error
}
