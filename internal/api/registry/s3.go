package registry

import (
	"github.com/mandarine-io/Backend/internal/api/config"
	"github.com/mandarine-io/Backend/pkg/storage/s3"
	s3Resource "github.com/mandarine-io/Backend/pkg/storage/s3/minio"
	"github.com/rs/zerolog/log"
)

func setupS3(c *Container) {
	log.Debug().Msg("setup s3")
	switch c.Config.S3.Type {
	case config.MinioS3Type:
		if c.Config.S3.Minio == nil {
			log.Fatal().Msg("minio config is nil")
		}
		minioConfig := mapAppMinioConfigToMinioConfig(&c.Config.S3)
		c.S3 = s3Resource.MustNewMinioClient(minioConfig)
		c.S3Client = s3.NewClient(c.S3, c.Config.S3.Minio.Bucket)
	default:
		log.Fatal().Msgf("unknown s3 type: %s", c.Config.S3.Type)
	}
}

func mapAppMinioConfigToMinioConfig(cfg *config.S3Config) *s3Resource.Config {
	return &s3Resource.Config{
		Address:    cfg.Minio.Address,
		AccessKey:  cfg.Minio.AccessKey,
		SecretKey:  cfg.Minio.SecretKey,
		BucketName: cfg.Minio.Bucket,
	}
}
