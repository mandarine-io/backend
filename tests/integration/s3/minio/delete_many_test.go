package minio

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/infrastructure/s3"
	"github.com/minio/minio-go/v7"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"golang.org/x/sys/unix"
	"os"
	"time"
)

var (
	deleteManyCh       chan string
	deleteManyDeleteCh chan string
)

type DeleteManySuite struct {
	suite.Suite
}

func (s *DeleteManySuite) BeforeEach(t provider.T) {
	t.Title("Delete Many - before each")
	t.Feature("Minio S3 manager")

	deleteManyCh = make(chan string, 10)
	deleteManyDeleteCh = make(chan string, 10)

	var (
		err   error
		files = make([]*os.File, 3)
	)

	for i := 0; i < 3; i++ {
		files[i], err = os.CreateTemp("", "delete-many.*.txt")
		t.Require().NoError(err)

		size, err := files[i].Write([]byte("content"))
		t.Require().NoError(err)

		_, err = files[i].Seek(0, unix.SEEK_SET)
		t.Require().NoError(err)

		info, err := client.PutObject(
			ctx, cfg.BucketName, uuid.New().String(), files[i], int64(size),
			minio.PutObjectOptions{
				SendContentMd5:        true,
				PartSize:              10 * 1024 * 1024,
				ConcurrentStreamParts: true,
				ContentType:           "text/plain",
				UserMetadata: map[string]string{
					s3.OriginalFilenameMetadata: files[i].Name(),
				},
			},
		)
		t.Require().NoError(err)

		deleteManyCh <- info.Key
		deleteManyDeleteCh <- info.Key
	}
}

func (s *DeleteManySuite) AfterEach(t provider.T) {
	t.Title("Delete Many - after each")
	t.Feature("Minio S3 manager")

	close(deleteManyCh)
	close(deleteManyDeleteCh)

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
	defer cancel()

	select {
	case objectId := <-deleteManyDeleteCh:
		err := client.RemoveObject(ctx, cfg.BucketName, objectId, minio.RemoveObjectOptions{})
		t.Require().NoError(err)
	case <-timeoutCtx.Done():
		return
	}
}

func (s *DeleteManySuite) Test_Success(t provider.T) {
	t.Title("Delete Many - success")
	t.Severity(allure.NORMAL)
	t.Feature("Minio S3 manager")
	t.Tags("Positive")

	objectIds := make([]string, 0, 10)

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
	defer cancel()

	select {
	case objectId := <-deleteManyCh:
		objectIds = append(objectIds, objectId)
	case <-timeoutCtx.Done():
		t.Require().NoError(timeoutCtx.Err())
	}

	errs := manager.DeleteMany(ctx, objectIds)

	for _, err := range errs {
		t.Require().NoError(err)
	}
}
