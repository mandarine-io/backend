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
	getOneCh       chan string
	getOneDeleteCh chan string
)

type GetOneSuite struct {
	suite.Suite
}

func (s *GetOneSuite) BeforeEach(t provider.T) {
	t.Title("Get One - before each")
	t.Feature("Minio S3 manager")

	getOneCh = make(chan string, 10)
	getOneDeleteCh = make(chan string, 10)

	file, err := os.CreateTemp("", "get-one.*.txt")
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

	getOneCh <- info.Key
	getOneDeleteCh <- info.Key
}

func (s *GetOneSuite) AfterEach(t provider.T) {
	t.Title("Get One - after each")
	t.Feature("Minio S3 manager")

	close(getOneCh)
	close(getOneDeleteCh)

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
	defer cancel()

	select {
	case objectId := <-getOneDeleteCh:
		err := client.RemoveObject(ctx, cfg.BucketName, objectId, minio.RemoveObjectOptions{})
		t.Require().NoError(err)
	case <-timeoutCtx.Done():
		return
	}
}

func (s *GetOneSuite) Test_Success(t provider.T) {
	t.Title("Get One - success")
	t.Severity(allure.NORMAL)
	t.Feature("Minio S3 manager")
	t.Tags("Positive")

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
	defer cancel()

	select {
	case objectId := <-getOneCh:
		res := manager.GetOne(ctx, objectId)
		t.Require().NoError(res.Error)
	case <-timeoutCtx.Done():
		t.Require().NoError(timeoutCtx.Err())
	}
}

func (s *GetOneSuite) Test_NotFound(t provider.T) {
	t.Title("Get One - not found")
	t.Severity(allure.CRITICAL)
	t.Feature("Minio S3 manager")
	t.Tags("Negative")

	res := manager.GetOne(ctx, "no-such-object")
	t.Require().Error(res.Error)
	t.Require().ErrorIs(res.Error, s3.ErrObjectNotFound)
}
