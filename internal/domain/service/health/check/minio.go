package check

import (
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

type MinioCheck struct {
	client *minio.Client
}

func NewMinioCheck(client *minio.Client) *MinioCheck {
	return &MinioCheck{client: client}
}

func (r *MinioCheck) Pass() bool {
	log.Debug().Msg("check minio connection")
	return r.client.IsOnline()
}

func (r *MinioCheck) Name() string {
	return "minio"
}
