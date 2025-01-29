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
	createOneCh chan string
)

type CreateOneSuite struct {
	suite.Suite
}

func (s *CreateOneSuite) BeforeEach(t provider.T) {
	t.Title("Create One - before each")
	t.Feature("Minio S3 manager")

	createOneCh = make(chan string, 10)
}
func (s *CreateOneSuite) AfterEach(t provider.T) {
	t.Title("Create One - after each")
	t.Feature("Minio S3 manager")

	close(createOneCh)

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
	defer cancel()

	select {
	case objectId := <-createOneCh:
		err := client.RemoveObject(ctx, cfg.BucketName, objectId, minio.RemoveObjectOptions{})
		t.Require().NoError(err)
	case <-timeoutCtx.Done():
		return
	}
}

func (s *CreateOneSuite) Test_Success(t provider.T) {
	t.Title("Create One - success")
	t.Severity(allure.NORMAL)
	t.Feature("Minio S3 manager")
	t.Tags("Positive")

	file, err := os.CreateTemp("", "create-one.*.txt")
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

	fileData := &s3.FileData{
		ID:          uuid.New().String(),
		Size:        int64(size),
		ContentType: "text/plain",
		Reader:      file,
		UserMetadata: map[string]string{
			s3.OriginalFilenameMetadata: file.Name(),
		},
	}

	res := manager.CreateOne(ctx, fileData)
	t.Require().NoError(res.Error)

	createOneCh <- res.ObjectID
}
