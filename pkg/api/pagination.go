package api

import (
	"encoding/base64"
	"errors"
)

// CursorPage holds keyset pagination results.
type CursorPage[T any] struct {
	Data []T        `json:"data"`
	Page CursorMeta `json:"page"`
}

// CursorMeta carries pagination metadata for cursor-based lists.
type CursorMeta struct {
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	Limit      int    `json:"limit"`
	HasMore    bool   `json:"has_more"`
}

// CursorPaginateOpts configures keyset pagination parameters.
type CursorPaginateOpts struct {
	Cursor string
	Limit  int
}

// Normalize ensures sane pagination defaults.
func (o *CursorPaginateOpts) Normalize() {
	if o.Limit <= 0 {
		o.Limit = 25
	}
	if o.Limit > 200 {
		o.Limit = 200
	}
}

// EncodeCursor serializes a cursor string for client URLs.
func EncodeCursor(raw string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// DecodeCursor decodes a cursor string from client input.
func DecodeCursor(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", errors.New("invalid cursor")
	}
	return string(raw), nil
}
