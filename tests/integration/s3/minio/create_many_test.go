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
	createManyCh chan string
)

type CreateManySuite struct {
	suite.Suite
}

func (s *CreateManySuite) BeforeEach(t provider.T) {
	t.Title("Create Many - before each")
	t.Feature("Minio S3 manager")

	createManyCh = make(chan string, 10)
}

func (s *CreateManySuite) AfterEach(t provider.T) {
	t.Title("Create Many - after each")
	t.Feature("Minio S3 manager")

	close(createManyCh)

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
	defer cancel()

	select {
	case objectId := <-createManyCh:
		err := client.RemoveObject(ctx, cfg.BucketName, objectId, minio.RemoveObjectOptions{})
		t.Require().NoError(err)
	case <-timeoutCtx.Done():
		return
	}
}

func (s *CreateManySuite) Test_Success(t provider.T) {
	t.Title("Create Many - success")
	t.Severity(allure.NORMAL)
	t.Feature("Minio S3 manager")
	t.Tags("Positive")

	var (
		err       error
		files     = make([]*os.File, 3)
		fileDatas = make([]*s3.FileData, 3)
	)

	for i := range fileDatas {
		files[i], err = os.CreateTemp("", "create-many.*.txt")
		t.Require().NoError(err)

		size, err := files[i].Write([]byte("content"))
		t.Require().NoError(err)

		_, err = files[i].Seek(0, unix.SEEK_SET)
		t.Require().NoError(err)

		fileDatas[i] = &s3.FileData{
			ID:          uuid.New().String(),
			Size:        int64(size),
			ContentType: "text/plain",
			Reader:      files[i],
			UserMetadata: map[string]string{
				s3.OriginalFilenameMetadata: files[i].Name(),
			},
		}
	}
	defer func(files []*os.File) {
		for _, file := range files {
			err := file.Close()
			t.Require().NoError(err)

			err = os.Remove(file.Name())
			t.Require().NoError(err)
		}
	}(files)

	results := manager.CreateMany(ctx, fileDatas)

	for _, result := range results {
		t.Require().NoError(result.Error)
		createManyCh <- result.ObjectID
	}
}
