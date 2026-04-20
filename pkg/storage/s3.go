package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// S3Driver handles S3-compatible storage (AWS S3, MinIO, Cloudflare R2, DigitalOcean Spaces).
type S3Driver struct {
	client       *s3.Client
	bucket       string
	publicURL    string
	usePathStyle bool
	downloader   *manager.Downloader
	uploader     *manager.Uploader
}

// NewS3Driver creates a new S3-compatible storage driver.
func NewS3Driver(region, bucket, endpoint, accessKey, secretKey, publicURL string, usePathStyle bool) (*S3Driver, error) {
	cfg := aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
	}

	// Use custom endpoint if provided (MinIO, R2, etc.)
	if endpoint != "" {
		cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		})
	}

	opts := []func(*s3.Options){
		func(o *s3.Options) {
			o.UsePathStyle = usePathStyle
		},
	}

	client := s3.NewFromConfig(cfg, opts...)

	d := &S3Driver{
		client:       client,
		bucket:       bucket,
		publicURL:    publicURL,
		usePathStyle: usePathStyle,
		downloader:   manager.NewDownloader(client),
		uploader:     manager.NewUploader(client),
	}

	return d, nil
}

// Put uploads a file to S3.
func (d *S3Driver) Put(ctx context.Context, path string, file any) error {
	path = d.cleanPath(path)

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

	_, err := d.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("s3 put: %w", err)
	}

	return nil
}

// Get downloads a file from S3.
func (d *S3Driver) Get(ctx context.Context, path string) ([]byte, error) {
	path = d.cleanPath(path)

	var buf manager.WriteAtBuffer
	_, err := d.downloader.Download(ctx, &buf, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NoSuchKey" {
				return nil, fmt.Errorf("file not found: %s", path)
			}
		}
		return nil, fmt.Errorf("s3 get: %w", err)
	}

	return buf.Bytes(), nil
}

// GetStream returns a readable stream for the file.
func (d *S3Driver) GetStream(ctx context.Context, path string) (io.ReadCloser, error) {
	path = d.cleanPath(path)

	result, err := d.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NoSuchKey" {
				return nil, fmt.Errorf("file not found: %s", path)
			}
		}
		return nil, fmt.Errorf("s3 get: %w", err)
	}

	return result.Body, nil
}

// Delete removes a file from S3.
func (d *S3Driver) Delete(ctx context.Context, path string) error {
	path = d.cleanPath(path)

	_, err := d.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("s3 delete: %w", err)
	}

	return nil
}

// Exists checks if a file exists on S3.
func (d *S3Driver) Exists(ctx context.Context, path string) (bool, error) {
	path = d.cleanPath(path)

	_, err := d.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NotFound" {
				return false, nil
			}
		}
		return false, fmt.Errorf("s3 head: %w", err)
	}

	return true, nil
}

// URL returns the public URL for the file.
func (d *S3Driver) URL(path string) string {
	path = d.cleanPath(path)

	if d.publicURL != "" {
		return d.publicURL + "/" + path
	}

	// Generate default S3 URL based on style
	if d.usePathStyle {
		return fmt.Sprintf("https://s3.amazonaws.com/%s/%s", d.bucket, path)
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", d.bucket, path)
}

// TemporaryURL generates a presigned/temporary URL valid for the given duration.
func (d *S3Driver) TemporaryURL(ctx context.Context, path string, duration time.Duration) (string, error) {
	path = d.cleanPath(path)

	presigner := s3.NewPresignClient(d.client)
	req, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})
	if err != nil {
		return "", fmt.Errorf("presign: %w", err)
	}

	return req.URL, nil
}

// Copy copies an object from source to destination.
func (d *S3Driver) Copy(ctx context.Context, sourcePath, destPath string) error {
	sourcePath = d.cleanPath(sourcePath)
	destPath = d.cleanPath(destPath)

	copySource := fmt.Sprintf("%s/%s", d.bucket, sourcePath)
	_, err := d.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(d.bucket),
		CopySource: aws.String(copySource),
		Key:        aws.String(destPath),
	})
	if err != nil {
		return fmt.Errorf("s3 copy: %w", err)
	}

	return nil
}

// Move moves (renames) an object.
func (d *S3Driver) Move(ctx context.Context, sourcePath, destPath string) error {
	// Copy then delete
	if err := d.Copy(ctx, sourcePath, destPath); err != nil {
		return err
	}
	return d.Delete(ctx, sourcePath)
}

// ListPrefix lists all objects with the given prefix.
func (d *S3Driver) ListPrefix(ctx context.Context, prefix string) ([]string, error) {
	prefix = d.cleanPath(prefix)
	if !strings.HasSuffix(prefix, "/") && prefix != "" {
		prefix += "/"
	}

	var keys []string
	paginator := s3.NewListObjectsV2Paginator(d.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(d.bucket),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("s3 list: %w", err)
		}

		for _, obj := range page.Contents {
			if obj.Key != nil {
				keys = append(keys, *obj.Key)
			}
		}
	}

	return keys, nil
}

// Size returns the file size in bytes.
func (d *S3Driver) Size(ctx context.Context, path string) (int64, error) {
	path = d.cleanPath(path)

	result, err := d.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NotFound" {
				return 0, fmt.Errorf("file not found: %s", path)
			}
		}
		return 0, fmt.Errorf("s3 head: %w", err)
	}

	return result.ContentLength, nil
}

// Name returns the driver name.
func (d *S3Driver) Name() string {
	return "s3"
}

// SetACL sets the object ACL (e.g., "public-read", "private").
func (d *S3Driver) SetACL(ctx context.Context, path, acl string) error {
	path = d.cleanPath(path)

	_, err := d.client.PutObjectAcl(ctx, &s3.PutObjectAclInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
		ACL:    types.ObjectCannedACL(acl),
	})
	if err != nil {
		return fmt.Errorf("s3 put acl: %w", err)
	}

	return nil
}

// Private helper methods

// cleanPath normalizes a file path (removes ../, ./, etc.)
func (d *S3Driver) cleanPath(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = filepath.Clean(path)
	path = strings.ReplaceAll(path, "\\", "/")
	// S3 doesn't use backslashes
	return path
}
