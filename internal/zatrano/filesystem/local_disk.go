package filesystem

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalDisk, yerel dosya sistemi için Disk arayüzünü uygular.
type LocalDisk struct {
	root string
	url  string
}

// NewLocalDisk, yerel disk sürücüsünü başlatır.
func NewLocalDisk(config DiskConfig) (Disk, error) {
	root, ok := config["root"]
	if !ok {
		return nil, fmt.Errorf("local disk requires a 'root' path")
	}
	// root klasörünün var olduğundan emin ol
	if err := os.MkdirAll(root, os.ModePerm); err != nil {
		return nil, fmt.Errorf("could not create root directory for local disk: %w", err)
	}
	return &LocalDisk{
		root: root,
		url:  config["url"],
	}, nil
}

// Put, bir dosyayı yerel diske kaydeder.
func (d *LocalDisk) Put(path string, fileHeader *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// `path`, `avatars` gibi bir alt klasör olabilir.
	// `filepath.Join` bunları güvenli bir şekilde birleştirir.
	relativePath := filepath.Join(path, fileName)
	fullPath := filepath.Join(d.root, relativePath)

	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		return "", err
	}
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	// İşletim sistemine özgü yolları (örn: `\`) URL uyumlu hale getir (`/`).
	return filepath.ToSlash(relativePath), nil
}

// --- YENİ EKLENEN METODLAR ---

// Get, yerel diskten bir dosyanın içeriğini byte dizisi olarak okur.
func (d *LocalDisk) Get(path string) ([]byte, error) {
	fullPath := filepath.Join(d.root, path)
	if !d.Exists(path) {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	return os.ReadFile(fullPath)
}

// Delete, yerel diskten bir dosyayı siler.
func (d *LocalDisk) Delete(path string) error {
	fullPath := filepath.Join(d.root, path)
	if !d.Exists(path) {
		// Dosya zaten yoksa, hata olarak kabul etmeyebiliriz.
		return nil
	}
	return os.Remove(fullPath)
}

// URL, bir dosyanın public URL'ini oluşturur.
func (d *LocalDisk) URL(path string) (string, error) {
	if d.url == "" {
		return "", fmt.Errorf("url not configured for this disk, it may not be a public disk")
	}
	// URL'lerin her zaman slash (/) kullanmasını sağla
	cleanPath := filepath.ToSlash(path)
	return strings.TrimSuffix(d.url, "/") + "/" + strings.TrimPrefix(cleanPath, "/"), nil
}

// Exists, bir dosyanın yerel diskte var olup olmadığını kontrol eder.
func (d *LocalDisk) Exists(path string) bool {
	fullPath := filepath.Join(d.root, path)
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
