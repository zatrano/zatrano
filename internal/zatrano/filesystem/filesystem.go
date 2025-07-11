package filesystem

import (
	"fmt"
	"mime/multipart"
)

// Filesystem, Storage arayüzünün somut uygulamasıdır.
type Filesystem struct {
	disks       map[string]Disk
	defaultDisk string
}

// NewFilesystem, yapılandırmaya göre yeni bir dosya sistemi yöneticisi oluşturur.
func NewFilesystem(config FilesystemConfig) (Storage, error) {
	fs := &Filesystem{
		disks:       make(map[string]Disk),
		defaultDisk: config.Default,
	}

	for name, diskConfig := range config.Disks {
		driver := diskConfig["driver"]
		var disk Disk
		var err error

		switch driver {
		case "local":
			disk, err = NewLocalDisk(diskConfig)
		// case "s3":
		// 	disk, err = NewS3Disk(diskConfig)
		default:
			return nil, fmt.Errorf("unsupported filesystem driver: %s", driver)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to initialize disk '%s': %w", name, err)
		}
		fs.disks[name] = disk
	}

	return fs, nil
}

// Disk, belirli bir diski seçer.
func (fs *Filesystem) Disk(name string) (Disk, error) {
	disk, ok := fs.disks[name]
	if !ok {
		return nil, fmt.Errorf("disk '%s' not found", name)
	}
	return disk, nil
}

// Varsayılan disk üzerinde çalışan helper metodlar
func (fs *Filesystem) Put(path string, file *multipart.FileHeader) (string, error) {
	d, err := fs.Disk(fs.defaultDisk)
	if err != nil {
		return "", err
	}
	return d.Put(path, file)
}

// Put, varsayılan disk üzerinde bir dosya kaydeder.
func (fs *Filesystem) Put(path string, file *multipart.FileHeader) (string, error) {
	d, err := fs.Disk(fs.defaultDisk)
	if err != nil {
		return "", err
	}
	return d.Put(path, file)
}

// Get, varsayılan diskten bir dosyanın içeriğini okur.
func (fs *Filesystem) Get(path string) ([]byte, error) {
	d, err := fs.Disk(fs.defaultDisk)
	if err != nil {
		return nil, err
	}
	return d.Get(path)
}

// Delete, varsayılan diskten bir dosyayı siler.
func (fs *Filesystem) Delete(path string) error {
	d, err := fs.Disk(fs.defaultDisk)
	if err != nil {
		return err
	}
	return d.Delete(path)
}

// URL, varsayılan diskteki bir dosyanın public URL'ini oluşturur.
func (fs *Filesystem) URL(path string) (string, error) {
	d, err := fs.Disk(fs.defaultDisk)
	if err != nil {
		return "", err
	}
	return d.URL(path)
}

// Exists, varsayılan diskte bir dosyanın var olup olmadığını kontrol eder.
func (fs *Filesystem) Exists(path string) bool {
	d, err := fs.Disk(fs.defaultDisk)
	if err != nil {
		return false
	}
	return d.Exists(path)
}
