package resource

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
	"log/slog"
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
		log.Fatal().Stack().Err(err).Msg("failed to connect to minio")
	}
	log.Info().Msgf("connected to minio host %s", endpoint)

	// Check if bucket exists
	log.Info().Msgf("check bucket \"%s\"", cfg.BucketName)
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to check minio bucket")
	}
	if !exists {
		slog.Info(fmt.Sprintf("Create bucket \"%s\"", cfg.BucketName))
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal().Stack().Err(err).Msg("failed to create minio bucket")
		}
	}

	return client
}
