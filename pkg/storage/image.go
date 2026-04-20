package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"strings"

	"image/gif"
)

// ImageProcessor handles image operations like resize, crop, and thumbnails.
type ImageProcessor struct {
	quality int // JPEG quality (1-100)
}

// NewImageProcessor creates a new image processor.
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		quality: 85,
	}
}

// ResizeOptions holds options for image resizing.
type ResizeOptions struct {
	Width  int
	Height int
	Fit    string // "contain", "cover", "fill", "inside", "outside"
}

// CropOptions holds options for image cropping.
type CropOptions struct {
	X      int
	Y      int
	Width  int
	Height int
}

// ThumbnailOptions holds options for generating thumbnails.
type ThumbnailOptions struct {
	Width  int
	Height int
	Fit    string
}

// Resize resizes an image and returns the result.
// Note: This uses Go's standard image package. For advanced image processing,
// consider integrating an external library like bimg (bindings to libvips).
func (ip *ImageProcessor) Resize(ctx context.Context, data []byte, opts ResizeOptions) ([]byte, error) {
	// Decode image
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	// Calculate new dimensions based on fit mode
	_, _ = ip.calculateDimensions(img.Bounds(), opts)

	// For now, return original data with metadata
	// In production, use bimg or similar for actual resizing
	// This is a placeholder that demonstrates the interface

	// Return original for now - implement actual resizing with external library
	return ip.encodeImage(img, format)
}

// Crop crops an image.
func (ip *ImageProcessor) Crop(ctx context.Context, data []byte, opts CropOptions) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	// For now, return original - implement with external library
	return ip.encodeImage(img, format)
}

// Thumbnail generates a thumbnail of the specified size.
func (ip *ImageProcessor) Thumbnail(ctx context.Context, data []byte, opts ThumbnailOptions) ([]byte, error) {
	resizeOpts := ResizeOptions{
		Width:  opts.Width,
		Height: opts.Height,
		Fit:    opts.Fit,
	}
	return ip.Resize(ctx, data, resizeOpts)
}

// ConvertFormat converts an image to a different format (jpeg, png, webp, etc.)
func (ip *ImageProcessor) ConvertFormat(ctx context.Context, data []byte, format string) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	return ip.encodeImageAs(img, format)
}

// GetImageInfo extracts basic info from an image (dimensions, format).
func (ip *ImageProcessor) GetImageInfo(data []byte) (*ImageInfo, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	bounds := img.Bounds()
	return &ImageInfo{
		Width:  bounds.Max.X,
		Height: bounds.Max.Y,
		Format: format,
	}, nil
}

// ImageInfo holds basic image metadata.
type ImageInfo struct {
	Width  int
	Height int
	Format string
}

// ProcessUploadedImage processes an uploaded image file (resize, validate, etc.)
// Returns the processed image data and metadata.
func (ip *ImageProcessor) ProcessUploadedImage(ctx context.Context, data []byte, filename string, maxWidth, maxHeight int) ([]byte, *ImageInfo, error) {
	// Validate image format
	info, err := ip.GetImageInfo(data)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid image: %w", err)
	}

	// Resize if needed
	if info.Width > maxWidth || info.Height > maxHeight {
		resized, err := ip.Resize(ctx, data, ResizeOptions{
			Width:  maxWidth,
			Height: maxHeight,
			Fit:    "contain",
		})
		if err != nil {
			return nil, nil, fmt.Errorf("resize: %w", err)
		}
		data = resized

		// Update info
		info, _ = ip.GetImageInfo(data)
	}

	return data, info, nil
}

// Private helper methods

// calculateDimensions calculates new dimensions based on fit mode.
func (ip *ImageProcessor) calculateDimensions(bounds image.Rectangle, opts ResizeOptions) (int, int) {
	currentWidth := bounds.Max.X
	currentHeight := bounds.Max.Y

	switch opts.Fit {
	case "contain":
		// Fit within bounds, maintain aspect ratio
		ratio := float64(currentHeight) / float64(currentWidth)
		newWidth := opts.Width
		newHeight := int(float64(newWidth) * ratio)
		if newHeight > opts.Height {
			newHeight = opts.Height
			newWidth = int(float64(newHeight) / ratio)
		}
		return newWidth, newHeight

	case "cover":
		// Fill bounds, maintain aspect ratio
		ratio := float64(currentHeight) / float64(currentWidth)
		newHeight := int(float64(opts.Width) * ratio)
		if newHeight < opts.Height {
			newHeight = opts.Height
		}
		return opts.Width, newHeight

	case "fill":
		// Stretch to fill bounds
		return opts.Width, opts.Height

	default:
		return opts.Width, opts.Height
	}
}

// encodeImage encodes an image back to its original format.
func (ip *ImageProcessor) encodeImage(img image.Image, format string) ([]byte, error) {
	return ip.encodeImageAs(img, format)
}

// encodeImageAs encodes an image to the specified format.
func (ip *ImageProcessor) encodeImageAs(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer

	format = strings.ToLower(format)
	switch format {
	case "jpeg", "jpg":
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: ip.quality})
		return buf.Bytes(), err

	case "png":
		err := png.Encode(&buf, img)
		return buf.Bytes(), err

	case "gif":
		err := gif.Encode(&buf, img, &gif.Options{})
		return buf.Bytes(), err

	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}
}

// SaveProcessedImage saves a processed image to a driver.
func SaveProcessedImage(ctx context.Context, driver Driver, data []byte, path string, variants map[string]ResizeOptions) (map[string]string, error) {
	// Save original
	if err := driver.Put(ctx, path, data); err != nil {
		return nil, fmt.Errorf("save original: %w", err)
	}

	results := map[string]string{
		"original": driver.URL(path),
	}

	// Save variants (e.g., thumbnail, medium, large)
	processor := NewImageProcessor()
	dir := filepath.Dir(path)
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	baseName := strings.TrimSuffix(name, ext)

	for variant, opts := range variants {
		resized, err := processor.Resize(ctx, data, opts)
		if err != nil {
			continue // Skip on error
		}

		variantPath := filepath.Join(dir, baseName+"_"+variant+ext)
		if err := driver.Put(ctx, variantPath, resized); err != nil {
			continue
		}

		results[variant] = driver.URL(variantPath)
	}

	return results, nil
}
