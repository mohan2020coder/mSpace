// internal/storage/minio.go
package storage

import (
	"context"

	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client     *minio.Client
	BucketName string
	Endpoint   string
	UseSSL     bool
}

func NewMinio(endpoint, accessKey, secretKey, bucket string, useSSL bool) *MinioClient {
	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("failed to init minio client: %v", err)
	}

	// Ensure bucket
	ctx := context.Background()
	exists, err := mc.BucketExists(ctx, bucket)
	if err != nil {
		// try to create
		if err := mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			log.Fatalf("failed to make bucket %s: %v", bucket, err)
		}
	} else if !exists {
		if err := mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			log.Fatalf("failed to make bucket %s: %v", bucket, err)
		}
	}

	return &MinioClient{
		Client:     mc,
		BucketName: bucket,
		Endpoint:   endpoint,
		UseSSL:     useSSL,
	}
}

// UploadStream puts an io.Reader to MinIO and returns object name
func (m *MinioClient) UploadStream(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := m.Client.PutObject(ctx, m.BucketName, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return objectName, nil
}

// PresignedURL returns a presigned GET URL for an object
// PresignedURL returns a presigned GET URL for an object
func (m *MinioClient) PresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := url.Values{} // ✅ correct type
	u, err := m.Client.PresignedGetObject(ctx, m.BucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
