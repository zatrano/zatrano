package storage

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"path/filepath"
	"testing"
	"time"
)

func TestLocalDriver_Put_Get(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	driver := NewLocalDriver(tmpDir, "/storage", true)
	ctx := context.Background()

	// Test data
	testFile := []byte("hello world")
	testPath := "documents/test.txt"

	// Put file
	err := driver.Put(ctx, testPath, testFile)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Get file
	data, err := driver.Get(ctx, testPath)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Verify
	if string(data) != string(testFile) {
		t.Errorf("Expected %q, got %q", testFile, data)
	}
}

func TestLocalDriver_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	driver := NewLocalDriver(tmpDir, "/storage", true)
	ctx := context.Background()

	// Non-existent file
	exists, err := driver.Exists(ctx, "nonexistent.txt")
	if err != nil || exists {
		t.Errorf("Expected false for non-existent file")
	}

	// Create file
	driver.Put(ctx, "exists.txt", "test")

	// Existent file
	exists, err = driver.Exists(ctx, "exists.txt")
	if err != nil || !exists {
		t.Errorf("Expected true for existing file")
	}
}

func TestLocalDriver_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	driver := NewLocalDriver(tmpDir, "/storage", true)
	ctx := context.Background()

	// Create file
	driver.Put(ctx, "delete-me.txt", "test")

	// Delete
	err := driver.Delete(ctx, "delete-me.txt")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	exists, _ := driver.Exists(ctx, "delete-me.txt")
	if exists {
		t.Errorf("File should be deleted")
	}
}

func TestLocalDriver_Copy(t *testing.T) {
	tmpDir := t.TempDir()
	driver := NewLocalDriver(tmpDir, "/storage", true)
	ctx := context.Background()

	// Create source
	driver.Put(ctx, "source.txt", "content")

	// Copy
	err := driver.Copy(ctx, "source.txt", "copied.txt")
	if err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	// Verify
	data, _ := driver.Get(ctx, "copied.txt")
	if string(data) != "content" {
		t.Errorf("Copy content mismatch")
	}
}

func TestLocalDriver_Move(t *testing.T) {
	tmpDir := t.TempDir()
	driver := NewLocalDriver(tmpDir, "/storage", true)
	ctx := context.Background()

	// Create file
	driver.Put(ctx, "original.txt", "test")

	// Move
	err := driver.Move(ctx, "original.txt", "moved.txt")
	if err != nil {
		t.Fatalf("Move failed: %v", err)
	}

	// Verify moved
	exists, _ := driver.Exists(ctx, "original.txt")
	if exists {
		t.Errorf("Original should be deleted after move")
	}

	exists, _ = driver.Exists(ctx, "moved.txt")
	if !exists {
		t.Errorf("Moved file should exist")
	}
}

func TestLocalDriver_Size(t *testing.T) {
	tmpDir := t.TempDir()
	driver := NewLocalDriver(tmpDir, "/storage", true)
	ctx := context.Background()

	// Create file
	testData := []byte("hello")
	driver.Put(ctx, "file.txt", testData)

	// Get size
	size, err := driver.Size(ctx, "file.txt")
	if err != nil {
		t.Fatalf("Size failed: %v", err)
	}

	if size != int64(len(testData)) {
		t.Errorf("Expected size %d, got %d", len(testData), size)
	}
}

func TestLocalDriver_URL(t *testing.T) {
	driver := NewLocalDriver("storage", "/storage", true)

	url := driver.URL("images/photo.jpg")
	if url != "/storage/images/photo.jpg" {
		t.Errorf("Expected /storage/images/photo.jpg, got %s", url)
	}
}

func TestLocalDriver_ListPrefix(t *testing.T) {
	tmpDir := t.TempDir()
	driver := NewLocalDriver(tmpDir, "/storage", true)
	ctx := context.Background()

	// Create files
	driver.Put(ctx, "dir/file1.txt", "content1")
	driver.Put(ctx, "dir/file2.txt", "content2")
	driver.Put(ctx, "dir/subdir/file3.txt", "content3")
	driver.Put(ctx, "other.txt", "other")

	// List dir
	files, err := driver.ListPrefix(ctx, "dir")
	if err != nil {
		t.Fatalf("ListPrefix failed: %v", err)
	}

	if len(files) < 3 {
		t.Errorf("Expected at least 3 files, got %d", len(files))
	}
}

func TestLocalDriver_TemporaryURL(t *testing.T) {
	driver := NewLocalDriver("storage", "/storage", true)
	ctx := context.Background()

	url, err := driver.TemporaryURL(ctx, "private/file.pdf", 15*time.Minute)
	if err != nil {
		t.Fatalf("TemporaryURL failed: %v", err)
	}

	if url == "" {
		t.Errorf("Expected non-empty URL")
	}

	// Should contain sig and expires params
	if !contains(url, "sig=") || !contains(url, "expires=") {
		t.Errorf("URL missing signature or expiration: %s", url)
	}
}

func TestImageProcessor_GetImageInfo(t *testing.T) {
	processor := NewImageProcessor()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}

	info, err := processor.GetImageInfo(buf.Bytes())
	if err != nil {
		t.Fatalf("GetImageInfo failed: %v", err)
	}

	if info.Width != 1 || info.Height != 1 {
		t.Errorf("Expected 1x1 image, got %dx%d", info.Width, info.Height)
	}

	if info.Format != "png" {
		t.Errorf("Expected format 'png', got %s", info.Format)
	}
}

func TestManager_Disk(t *testing.T) {
	manager := NewManager()

	// Register drivers
	local := NewLocalDriver("storage", "/storage", true)
	manager.Register("local", local)
	manager.SetDefault("local")

	// Get default
	driver, err := manager.Disk("")
	if err != nil {
		t.Fatalf("Failed to get default disk: %v", err)
	}

	if driver.Name() != "local" {
		t.Errorf("Expected local driver, got %s", driver.Name())
	}

	// Get by name
	driver, err = manager.Disk("local")
	if err != nil {
		t.Fatalf("Failed to get local disk: %v", err)
	}

	// Non-existent
	_, err = manager.Disk("nonexistent")
	if err == nil {
		t.Errorf("Expected error for non-existent disk")
	}
}

// Helper functions

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests

func BenchmarkLocalDriver_Put(b *testing.B) {
	tmpDir := b.TempDir()
	driver := NewLocalDriver(tmpDir, "/storage", true)
	ctx := context.Background()
	testData := []byte("test data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := filepath.Join("bench", "file"+string(rune(i))+".txt")
		driver.Put(ctx, path, testData)
	}
}

func BenchmarkLocalDriver_Get(b *testing.B) {
	tmpDir := b.TempDir()
	driver := NewLocalDriver(tmpDir, "/storage", true)
	ctx := context.Background()

	// Pre-populate files
	driver.Put(ctx, "bench/file.txt", []byte("test data"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		driver.Get(ctx, "bench/file.txt")
	}
}
