package dto

import "io"

type FileData struct {
	ID           string
	Size         int64
	ContentType  string
	Reader       io.ReadCloser
	UserMetadata map[string]string
}

type CreateDto struct {
	ObjectID string
	Error    error
}

type GetDto struct {
	Data  *FileData
	Error error
}
