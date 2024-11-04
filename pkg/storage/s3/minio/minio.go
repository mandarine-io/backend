package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
	"log/slog"
)

type Config struct {
	Address    string
	AccessKey  string
	SecretKey  string
	BucketName string
}

func MustNewMinioClient(cfg *Config) *minio.Client {
	// Configure to use MinIO Server
	ctx := context.Background()
	client, err := minio.New(cfg.Address, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to connect to minio")
	}
	log.Info().Msgf("connected to minio host %s", cfg.Address)

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
