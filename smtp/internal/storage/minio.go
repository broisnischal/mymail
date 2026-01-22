package storage

import (
	"context"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mymail/smtp/internal/config"
)

type MinIO struct {
	client *minio.Client
	bucket string
}

func NewMinIO(cfg config.MinIOConfig) (*MinIO, error) {
	endpoint := cfg.Endpoint
	useSSL := cfg.UseSSL

	// Parse endpoint
	parts := strings.Split(endpoint, ":")
	host := parts[0]
	port := "9000"
	if len(parts) > 1 {
		port = parts[1]
	}

	client, err := minio.New(host+":"+port, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinIO{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

func (m *MinIO) Upload(ctx context.Context, path string, reader io.Reader, size int64) error {
	_, err := m.client.PutObject(ctx, m.bucket, path, reader, size, minio.PutObjectOptions{
		ContentType: "message/rfc822",
	})
	return err
}

// UploadStream streams data directly to MinIO without requiring a known size upfront
// This is used for streaming email uploads where we don't know the total size
func (m *MinIO) UploadStream(ctx context.Context, path string, reader io.Reader) error {
	// Use -1 for size to indicate unknown size, MinIO will handle streaming
	_, err := m.client.PutObject(ctx, m.bucket, path, reader, -1, minio.PutObjectOptions{
		ContentType: "message/rfc822",
	})
	return err
}

func (m *MinIO) Get(ctx context.Context, path string) (io.Reader, error) {
	obj, err := m.client.GetObject(ctx, m.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *MinIO) Delete(ctx context.Context, path string) error {
	return m.client.RemoveObject(ctx, m.bucket, path, minio.RemoveObjectOptions{})
}
