package resource

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log/slog"
	"mandarine/pkg/logging"
	"os"
)

type MinioConfig struct {
	Host       string
	Port       int
	AccessKey  string
	SecretKey  string
	BucketName string
}

func MustConnectMinio(cfg *MinioConfig) *minio.Client {
	// Configure to use MinIO Server
	ctx := context.Background()
	endpoint := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		slog.Error("MinIO connection error", logging.ErrorAttr(err))
		os.Exit(1)
	}

	// Check if bucket exists
	slog.Info(fmt.Sprintf("Check bucket \"%s\"", cfg.BucketName))
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		slog.Error("MinIO bucket checking error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	if !exists {
		slog.Info(fmt.Sprintf("Create bucket \"%s\"", cfg.BucketName))
		err := client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			slog.Error("MinIO bucket creation error", logging.ErrorAttr(err))
			os.Exit(1)
		}
	}

	slog.Info("Connected to Minio host " + endpoint)

	return client
}
