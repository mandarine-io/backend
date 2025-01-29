package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

type Config struct {
	Address    string
	AccessKey  string
	SecretKey  string
	BucketName string
}

func NewClient(cfg Config) (*minio.Client, error) {
	client, err := minio.New(
		cfg.Address, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: false,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to client: %w", err)
	}

	// Check if bucket exists
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check client bucket: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create client bucket: %w", err)
		}
	}

	return client, nil
}
