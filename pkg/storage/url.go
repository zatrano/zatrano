package storage

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// URLHelper provides utility functions for working with storage URLs.
type URLHelper struct {
	drivers map[string]Driver
}

// NewURLHelper creates a new URL helper.
func NewURLHelper(drivers map[string]Driver) *URLHelper {
	return &URLHelper{
		drivers: drivers,
	}
}

// TemporaryURLOptions holds options for generating temporary URLs.
type TemporaryURLOptions struct {
	Driver   string            // Which disk to use (default: first registered)
	Path     string            // File path
	Duration time.Duration     // How long the URL is valid
	Query    map[string]string // Additional query parameters
}

// TemporaryURL generates a temporary signed URL for a file.
func (h *URLHelper) TemporaryURL(ctx context.Context, opts TemporaryURLOptions) (string, error) {
	driver := h.getDriver(opts.Driver)
	if driver == nil {
		return "", fmt.Errorf("driver not found: %s", opts.Driver)
	}

	if opts.Duration == 0 {
		opts.Duration = 15 * time.Minute // Default 15 minutes
	}

	tempURL, err := driver.TemporaryURL(ctx, opts.Path, opts.Duration)
	if err != nil {
		return "", err
	}

	// Add additional query parameters if provided
	if len(opts.Query) > 0 {
		parsed, _ := url.Parse(tempURL)
		q := parsed.Query()
		for k, v := range opts.Query {
			q.Set(k, v)
		}
		parsed.RawQuery = q.Encode()
		tempURL = parsed.String()
	}

	return tempURL, nil
}

// AssetURL returns a URL for a public asset (with cache-busting hash if applicable).
// This is different from temporary URLs - it's for permanent public files.
func (h *URLHelper) AssetURL(driver, path string) string {
	d := h.getDriver(driver)
	if d == nil {
		return ""
	}
	return d.URL(path)
}

// TemporaryDownloadURL generates a temporary URL for downloading a file.
// It includes a Content-Disposition header suggestion (filename).
func (h *URLHelper) TemporaryDownloadURL(ctx context.Context, opts TemporaryURLOptions) (string, error) {
	if opts.Query == nil {
		opts.Query = make(map[string]string)
	}

	// Extract filename from path
	filename := extractFilename(opts.Path)
	opts.Query["Content-Disposition"] = "attachment; filename=\"" + filename + "\""

	return h.TemporaryURL(ctx, opts)
}

// TemporaryPreviewURL generates a temporary URL for previewing a file (inline).
func (h *URLHelper) TemporaryPreviewURL(ctx context.Context, opts TemporaryURLOptions) (string, error) {
	if opts.Query == nil {
		opts.Query = make(map[string]string)
	}
	opts.Query["Content-Disposition"] = "inline"

	return h.TemporaryURL(ctx, opts)
}

// URLWithDuration returns a temporary URL with a custom duration.
// Shorthand for TemporaryURL.
func (h *URLHelper) URLWithDuration(ctx context.Context, driver, path string, duration time.Duration) (string, error) {
	return h.TemporaryURL(ctx, TemporaryURLOptions{
		Driver:   driver,
		Path:     path,
		Duration: duration,
	})
}

// URLValid returns true if a temporary URL is still valid (before expiration).
func URLValid(urlStr string) bool {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	expiresStr := parsed.Query().Get("expires")
	if expiresStr == "" {
		return true // No expiration
	}

	expires, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil {
		return false
	}

	return time.Now().Unix() < expires
}

// ExtractPath extracts the storage path from a URL.
// For example, "/storage/images/photo.jpg" -> "images/photo.jpg"
func ExtractPath(urlStr, prefix string) string {
	if prefix == "" {
		prefix = "/storage/"
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	path := parsed.Path
	if len(path) > len(prefix) && path[:len(prefix)] == prefix {
		return path[len(prefix):]
	}

	return path
}

// Private helper methods

// getDriver returns a driver from the map, or the first one if not specified.
func (h *URLHelper) getDriver(name string) Driver {
	if name == "" {
		// Return first registered driver
		for _, d := range h.drivers {
			return d
		}
		return nil
	}
	return h.drivers[name]
}

// extractFilename extracts the filename from a path.
func extractFilename(path string) string {
	// Find last slash
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}
	return path
}
