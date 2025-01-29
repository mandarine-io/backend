package initializer

import (
	"github.com/mandarine-io/backend/config"
	"github.com/mandarine-io/backend/internal/di"
	minio2 "github.com/mandarine-io/backend/internal/infrastructure/s3/minio"
)

func S3(c *di.Container) di.Initializer {
	return func() error {
		c.Logger.Debug().Msg("setup s3")

		var err error
		minioCfg := toMinioS3Config(c.Config.S3)

		c.Infrastructure.MinioClient, err = minio2.NewClient(minioCfg)
		if err != nil {
			return err
		}

		c.Logger.Info().Msgf("connect to minio %s", minioCfg.Address)

		c.Infrastructure.S3Manager, err = minio2.NewManager(
			c.Infrastructure.MinioClient,
			minioCfg.BucketName,
			minio2.WithLogger(c.Logger.With().Str("component", "minio-s3").Logger()),
		)

		return err
	}
}

func toMinioS3Config(cfg config.MinIOS3Config) minio2.Config {
	return minio2.Config{
		Address:    cfg.Address,
		AccessKey:  cfg.AccessKey,
		SecretKey:  cfg.SecretKey,
		BucketName: cfg.Bucket,
	}
}
