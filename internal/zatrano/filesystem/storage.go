package filesystem

import "mime/multipart"

// Storage, dosya işlemleri için ana arayüzdür.
type Storage interface {
	// Disk, işlemleri belirli bir disk üzerinde yapmak için seçer.
	Disk(name string) (Disk, error)

	// Varsayılan disk üzerinde dosya işlemleri
	Put(path string, file *multipart.FileHeader) (string, error)
	Get(path string) ([]byte, error)
	Delete(path string) error
	URL(path string) (string, error)
	Exists(path string) bool
}

// Disk, belirli bir depolama sürücüsü için işlemleri tanımlar.
type Disk interface {
	Put(path string, file *multipart.FileHeader) (string, error)
	Get(path string) ([]byte, error)
	Delete(path string) error
	URL(path string) (string, error)
	Exists(path string) bool
}

// DiskConfig, bir diskin yapılandırma seçeneklerini tutar.
type DiskConfig map[string]string
