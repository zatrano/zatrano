package storage

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// ServeFileMiddleware handles serving private files with temporary URL validation.
// Usage: app.Get("/storage/private/:path", storage.ServeFileMiddleware(storageManager))
func ServeFileMiddleware(manager *Manager) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		path := c.Params("path", "")
		if path == "" {
			return c.SendStatus(fiber.StatusNotFound)
		}

		// Get query parameters
		sig := c.Query("sig")
		expiresStr := c.Query("expires")

		// Validate signature if provided
		if sig != "" && expiresStr != "" {
			expires, err := strconv.ParseInt(expiresStr, 10, 64)
			if err != nil {
				return c.SendStatus(fiber.StatusUnauthorized)
			}

			if !VerifyTemporaryURL(path, sig, expires) {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
		} else {
			// No signature provided - deny access to private files
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Get private driver
		driver, err := manager.Disk("private")
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Retrieve file
		ctx := context.Background()
		stream, err := driver.GetStream(ctx, path)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		defer stream.Close()

		// Get file size
		size, _ := driver.Size(ctx, path)

		// Set response headers
		c.Set("Content-Length", strconv.FormatInt(size, 10))
		c.Set("Content-Type", getContentType(path))

		// Check Content-Disposition query param
		if disposition := c.Query("Content-Disposition"); disposition != "" {
			c.Set("Content-Disposition", disposition)
		}

		// Stream response
		return c.SendStream(stream)
	}
}

// UploadMiddleware handles file uploads with validation.
// Usage: app.Post("/upload", storage.UploadMiddleware(storageManager, uploadConfig))
func UploadMiddleware(manager *Manager, config *UploadConfig) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "file required",
			})
		}

		// Validate file size
		if config.MaxSize > 0 && file.Size > config.MaxSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "file too large",
			})
		}

		// Validate file type
		if len(config.AllowedMimeTypes) > 0 {
			if !isAllowedMimeType(file.Header.Get("Content-Type"), config.AllowedMimeTypes) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "file type not allowed",
				})
			}
		}

		// Open file
		f, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to read file",
			})
		}
		defer f.Close()

		// Get driver
		diskName := c.FormValue("disk", "local")
		driver, err := manager.Disk(diskName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "storage driver not found",
			})
		}

		// Generate path
		path := config.PathPrefix + "/" + file.Filename
		if config.PathGenerator != nil {
			path = config.PathGenerator(file.Filename, file.Header.Get("Content-Type"))
		}

		// Store file
		ctx := c.Context()
		err = driver.Put(ctx, path, f)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to store file",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"file": fiber.Map{
				"name": file.Filename,
				"size": file.Size,
				"url":  driver.URL(path),
				"path": path,
			},
		})
	}
}

// UploadConfig holds configuration for file uploads.
type UploadConfig struct {
	MaxSize          int64
	AllowedMimeTypes []string
	PathPrefix       string
	PathGenerator    func(filename, mimeType string) string
	ValidateCallback func(filename string, mimeType string) error
}

// StaticFilesMiddleware serves files from storage like static files.
// Usage: app.Static("/storage", storageManager.Disk("public"))
func StaticFilesMiddleware(driver Driver, prefix string) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		path := c.Params("*1", "")
		if path == "" {
			return c.SendStatus(fiber.StatusNotFound)
		}

		ctx := c.Context()
		exists, _ := driver.Exists(ctx, path)
		if !exists {
			return c.SendStatus(fiber.StatusNotFound)
		}

		stream, err := driver.GetStream(ctx, path)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		defer stream.Close()

		// Set cache headers
		c.Set("Cache-Control", "public, max-age=31536000") // 1 year
		c.Set("Content-Type", getContentType(path))

		return c.SendStream(stream)
	}
}

// DownloadHandler provides a simple file download endpoint.
func DownloadHandler(manager *Manager) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		path := c.Query("path")
		disk := c.Query("disk", "local")

		if path == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		driver, err := manager.Disk(disk)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		ctx := c.Context()
		exists, _ := driver.Exists(ctx, path)
		if !exists {
			return c.SendStatus(fiber.StatusNotFound)
		}

		// Generate temporary URL if private
		if disk == "private" {
			tempURL, err := driver.TemporaryURL(ctx, path, 24*3600) // 1 hour
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			return c.Redirect().To(tempURL)
		}

		// For public files, stream directly
		stream, err := driver.GetStream(ctx, path)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		defer stream.Close()

		filename := extractFilename(path)
		c.Set("Content-Disposition", `attachment; filename="`+filename+`"`)
		c.Set("Content-Type", getContentType(path))

		return c.SendStream(stream)
	}
}

// Private helper functions

// getContentType returns the MIME type for a file based on its extension.
func getContentType(path string) string {
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".pdf":  "application/pdf",
		".txt":  "text/plain",
		".html": "text/html",
		".csv":  "text/csv",
		".json": "application/json",
		".zip":  "application/zip",
		".mp4":  "video/mp4",
		".mp3":  "audio/mpeg",
	}

	// Extract extension
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			ext := strings.ToLower(path[i:])
			if mt, ok := mimeTypes[ext]; ok {
				return mt
			}
			break
		}
	}

	return "application/octet-stream"
}

// isAllowedMimeType checks if a MIME type is in the allowed list.
func isAllowedMimeType(mimeType string, allowed []string) bool {
	// Extract main type (e.g., "image" from "image/jpeg")
	parts := strings.Split(mimeType, "/")
	if len(parts) < 1 {
		return false
	}

	mainType := parts[0]

	for _, a := range allowed {
		if a == mimeType || a == mainType || a == mainType+"/*" {
			return true
		}
	}

	return false
}

// ParseDownloadURL extracts storage parameters from a download URL.
func ParseDownloadURL(urlStr string) (disk, path string, err error) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	disk = parsed.Query().Get("disk")
	path = parsed.Query().Get("path")

	if path == "" {
		path = strings.TrimPrefix(parsed.Path, "/download/")
	}

	return
}
