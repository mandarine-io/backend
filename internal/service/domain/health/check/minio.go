package check

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
)

type MinioCheck struct {
	client *minio.Client
	logger zerolog.Logger
}

type MinioOption func(c *MinioCheck) error

func WithMinioLogger(logger zerolog.Logger) MinioOption {
	return func(c *MinioCheck) error {
		c.logger = logger
		return nil
	}
}

func NewMinioCheck(client *minio.Client, opts ...MinioOption) (*MinioCheck, error) {
	check := &MinioCheck{
		client: client,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		if err := opt(check); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return check, nil
}

func (c *MinioCheck) Pass() bool {
	c.logger.Debug().Msg("check minio connection")
	return c.client.IsOnline()
}

func (c *MinioCheck) Name() string {
	return "minio"
}
