package storage

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalDriver stores files on the local filesystem.
type LocalDriver struct {
	root      string // Root directory for storage
	publicURL string // Public URL prefix for public disks
	public    bool   // Whether this is a public disk
}

// NewLocalDriver creates a new local storage driver.
func NewLocalDriver(root, publicURL string, public bool) *LocalDriver {
	return &LocalDriver{
		root:      root,
		publicURL: publicURL,
		public:    public,
	}
}

// Put stores a file at the given path.
func (d *LocalDriver) Put(ctx context.Context, path string, file any) error {
	path = d.cleanPath(path)
	fullPath := filepath.Join(d.root, path)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	// Convert file to []byte
	var data []byte
	switch v := file.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	case io.Reader:
		b, err := io.ReadAll(v)
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}
		data = b
	default:
		return fmt.Errorf("unsupported file type: %T", file)
	}

	// Write file
	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

// Get retrieves the file contents.
func (d *LocalDriver) Get(ctx context.Context, path string) ([]byte, error) {
	path = d.cleanPath(path)
	fullPath := filepath.Join(d.root, path)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("read: %w", err)
	}

	return data, nil
}

// GetStream returns a readable stream for the file.
func (d *LocalDriver) GetStream(ctx context.Context, path string) (io.ReadCloser, error) {
	path = d.cleanPath(path)
	fullPath := filepath.Join(d.root, path)

	f, err := os.Open(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("open: %w", err)
	}

	return f, nil
}

// Delete removes the file.
func (d *LocalDriver) Delete(ctx context.Context, path string) error {
	path = d.cleanPath(path)
	fullPath := filepath.Join(d.root, path)

	if err := os.Remove(fullPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // Deletion of non-existent file is a no-op
		}
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Exists checks if the file exists.
func (d *LocalDriver) Exists(ctx context.Context, path string) (bool, error) {
	path = d.cleanPath(path)
	fullPath := filepath.Join(d.root, path)

	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// URL returns a public URL for the file.
func (d *LocalDriver) URL(path string) string {
	path = d.cleanPath(path)
	if d.publicURL == "" {
		return "/" + strings.ReplaceAll(path, "\\", "/")
	}
	return d.publicURL + "/" + strings.ReplaceAll(path, "\\", "/")
}

// TemporaryURL returns a temporary signed URL (local driver uses a simple hash).
// For local storage, this just returns the regular URL with a hash parameter.
func (d *LocalDriver) TemporaryURL(ctx context.Context, path string, duration time.Duration) (string, error) {
	path = d.cleanPath(path)

	// Generate a simple hash based on path and expiration time
	expiry := time.Now().Add(duration).Unix()
	hash := d.signPath(path, expiry)

	url := d.URL(path) + "?sig=" + hash + "&expires=" + fmt.Sprintf("%d", expiry)
	return url, nil
}

// Copy copies a file from source to destination.
func (d *LocalDriver) Copy(ctx context.Context, sourcePath, destPath string) error {
	sourcePath = d.cleanPath(sourcePath)
	destPath = d.cleanPath(destPath)

	srcFull := filepath.Join(d.root, sourcePath)
	dstFull := filepath.Join(d.root, destPath)

	// Read source
	data, err := os.ReadFile(srcFull)
	if err != nil {
		return fmt.Errorf("read source: %w", err)
	}

	// Create dest directory
	if err := os.MkdirAll(filepath.Dir(dstFull), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	// Write destination
	if err := os.WriteFile(dstFull, data, 0o644); err != nil {
		return fmt.Errorf("write destination: %w", err)
	}

	return nil
}

// Move renames/moves a file.
func (d *LocalDriver) Move(ctx context.Context, sourcePath, destPath string) error {
	sourcePath = d.cleanPath(sourcePath)
	destPath = d.cleanPath(destPath)

	srcFull := filepath.Join(d.root, sourcePath)
	dstFull := filepath.Join(d.root, destPath)

	// Create dest directory
	if err := os.MkdirAll(filepath.Dir(dstFull), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	// Rename file
	if err := os.Rename(srcFull, dstFull); err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	return nil
}

// ListPrefix lists all files under a prefix (directory).
func (d *LocalDriver) ListPrefix(ctx context.Context, prefix string) ([]string, error) {
	prefix = d.cleanPath(prefix)
	fullPath := filepath.Join(d.root, prefix)

	var files []string
	err := filepath.WalkDir(fullPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			// Get path relative to root
			rel, _ := filepath.Rel(d.root, path)
			files = append(files, strings.ReplaceAll(rel, "\\", "/"))
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil // Empty list if prefix doesn't exist
		}
		return nil, fmt.Errorf("walk: %w", err)
	}

	return files, nil
}

// Size returns the file size in bytes.
func (d *LocalDriver) Size(ctx context.Context, path string) (int64, error) {
	path = d.cleanPath(path)
	fullPath := filepath.Join(d.root, path)

	fi, err := os.Stat(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, fmt.Errorf("file not found: %s", path)
		}
		return 0, fmt.Errorf("stat: %w", err)
	}

	return fi.Size(), nil
}

// Name returns the driver name.
func (d *LocalDriver) Name() string {
	return "local"
}

// Private helper methods

// cleanPath normalizes a file path (removes ../, ./, etc.)
func (d *LocalDriver) cleanPath(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = filepath.Clean(path)
	path = strings.ReplaceAll(path, "\\", "/")
	return path
}

// signPath creates a simple HMAC-like signature for temporary URLs.
func (d *LocalDriver) signPath(path string, expiry int64) string {
	message := fmt.Sprintf("%s:%d", path, expiry)
	hash := md5.Sum([]byte(message))
	return hex.EncodeToString(hash[:])
}

// VerifyTemporaryURL checks if a temporary URL signature is valid (for middleware).
func VerifyTemporaryURL(path, sig string, expiry int64) bool {
	if time.Now().Unix() > expiry {
		return false
	}
	message := fmt.Sprintf("%s:%d", path, expiry)
	hash := md5.Sum([]byte(message))
	expectedSig := hex.EncodeToString(hash[:])
	return sig == expectedSig
}
