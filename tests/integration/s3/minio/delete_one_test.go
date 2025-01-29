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
	deleteOneCh       chan string
	deleteOneDeleteCh chan string
)

type DeleteOneSuite struct {
	suite.Suite
}

func (s *DeleteOneSuite) BeforeEach(t provider.T) {
	t.Title("Delete One - before each")
	t.Feature("Minio S3 manager")

	deleteOneCh = make(chan string, 10)
	deleteOneDeleteCh = make(chan string, 10)

	file, err := os.CreateTemp("", "delete-one.*.txt")
	t.Require().NoError(err)
	defer func(file *os.File) {
		err := file.Close()
		t.Require().NoError(err)

		err = os.Remove(file.Name())
		t.Require().NoError(err)
	}(file)

	size, err := file.Write([]byte("content"))
	t.Require().NoError(err)

	_, err = file.Seek(0, unix.SEEK_SET)
	t.Require().NoError(err)

	info, err := client.PutObject(
		ctx, cfg.BucketName, uuid.New().String(), file, int64(size),
		minio.PutObjectOptions{
			SendContentMd5:        true,
			PartSize:              10 * 1024 * 1024,
			ConcurrentStreamParts: true,
			ContentType:           "text/plain",
			UserMetadata: map[string]string{
				s3.OriginalFilenameMetadata: file.Name(),
			},
		},
	)
	t.Require().NoError(err)

	deleteOneCh <- info.Key
	deleteOneDeleteCh <- info.Key
}

func (s *DeleteOneSuite) AfterEach(t provider.T) {
	t.Title("Delete One - after each")
	t.Feature("Minio S3 manager")

	close(deleteOneCh)
	close(deleteOneDeleteCh)

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
	defer cancel()

	select {
	case objectId := <-deleteOneDeleteCh:
		err := client.RemoveObject(ctx, cfg.BucketName, objectId, minio.RemoveObjectOptions{})
		t.Require().NoError(err)
	case <-timeoutCtx.Done():
		return
	}
}

func (s *DeleteOneSuite) Test_Success(t provider.T) {
	t.Title("Delete One - success")
	t.Severity(allure.NORMAL)
	t.Feature("Minio S3 manager")
	t.Tags("Positive")

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
	defer cancel()

	select {
	case objectId := <-deleteOneCh:
		err := manager.DeleteOne(ctx, objectId)
		t.Require().NoError(err)
	case <-timeoutCtx.Done():
		t.Require().NoError(timeoutCtx.Err())
	}
}

func (s *DeleteOneSuite) Test_NotFound(t provider.T) {
	t.Title("Delete One - not found")
	t.Severity(allure.CRITICAL)
	t.Feature("Minio S3 manager")
	t.Tags("Negative")

	err := manager.DeleteOne(ctx, "no-such-object")
	t.Require().NoError(err)
}
