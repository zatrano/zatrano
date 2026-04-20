// Package storage provides a unified interface for file storage operations
// across different backends (local disk, S3, MinIO, Cloudflare R2, etc.).
package storage

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Driver is the interface that all storage drivers must implement.
type Driver interface {
	// Put stores a file at the given path. The file parameter can be a []byte,
	// io.Reader, or string.
	Put(ctx context.Context, path string, file any) error

	// Get retrieves the file contents from the given path.
	Get(ctx context.Context, path string) ([]byte, error)

	// GetStream returns an io.ReadCloser for streaming large files.
	GetStream(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete removes the file at the given path.
	Delete(ctx context.Context, path string) error

	// Exists checks whether a file exists at the given path.
	Exists(ctx context.Context, path string) (bool, error)

	// URL returns a public URL for the given path.
	// For local disk, this is a relative URL. For cloud drivers, it's the full URL.
	URL(path string) string

	// TemporaryURL returns a signed/temporary URL valid for the given duration.
	// Used for private files or time-limited access.
	TemporaryURL(ctx context.Context, path string, duration time.Duration) (string, error)

	// Copy copies a file from sourcePath to destPath on the same driver.
	Copy(ctx context.Context, sourcePath, destPath string) error

	// Move moves (renames) a file from sourcePath to destPath on the same driver.
	Move(ctx context.Context, sourcePath, destPath string) error

	// ListPrefix lists all files under the given prefix/directory.
	ListPrefix(ctx context.Context, prefix string) ([]string, error)

	// Size returns the file size in bytes.
	Size(ctx context.Context, path string) (int64, error)

	// Name returns the name of the driver (e.g., "local", "s3", "minio").
	Name() string
}

// Config holds common storage configuration.
type Config struct {
	// Driver type: "local", "s3", "minio", "r2"
	Driver string

	// LocalDisk config
	Local struct {
		// Root directory for local storage (default: "storage/app")
		Root string
		// PublicRoot directory for public files (default: "public/storage")
		PublicRoot string
		// PublicURL prefix for accessing public files (default: "/storage")
		PublicURL string
	}

	// S3/MinIO/R2 config
	S3 struct {
		// Region for AWS S3 or compatible service
		Region string
		// Bucket name
		Bucket string
		// Endpoint URL (for MinIO, Cloudflare R2, etc.; leave empty for AWS S3)
		Endpoint string
		// Access key ID
		AccessKey string
		// Secret access key
		SecretKey string
		// Public URL prefix (optional; overrides default bucket URL)
		PublicURL string
		// UsePathStyle uses path-style URLs (s3.com/bucket/key instead of bucket.s3.com/key)
		UsePathStyle bool
	}

	// Private disks configuration (e.g., for "private" named disk)
	Private struct {
		Enabled bool
		Root    string
	}
}

// Manager manages multiple named storage disks.
type Manager struct {
	drivers     map[string]Driver
	defaultDisk string
}

// NewManager creates a new storage manager.
func NewManager() *Manager {
	return &Manager{
		drivers: make(map[string]Driver),
	}
}

// Register registers a driver with the given name.
func (m *Manager) Register(name string, driver Driver) {
	m.drivers[name] = driver
	if m.defaultDisk == "" {
		m.defaultDisk = name
	}
}

// SetDefault sets the default driver.
func (m *Manager) SetDefault(name string) {
	m.defaultDisk = name
}

// Disk returns a named driver, or the default if name is empty.
func (m *Manager) Disk(name string) (Driver, error) {
	if name == "" {
		name = m.defaultDisk
	}
	d, ok := m.drivers[name]
	if !ok {
		return nil, fmt.Errorf("storage driver %q not found", name)
	}
	return d, nil
}

// Default returns the default driver.
func (m *Manager) Default() (Driver, error) {
	return m.Disk("")
}
