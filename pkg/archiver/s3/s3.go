package s3

import (
	"bytes"
	"context"
	"io"
)

// "github.com/aws/aws-sdk-go/service/s3/s3manager"
// https://godoc.org/github.com/aws/aws-sdk-go/service/s3/s3manager#NewUploader

type uploader interface {
	UploadWithContext(ctx context.Context, filename string, body io.Reader) (string, error)
}

// S3Archiver implements the Archiver interface, and archives to S3.
type S3Archiver struct {
	uploader uploader
}

// New creates and returns a new S3Archiver.
func New(u uploader) *S3Archiver {
	return &S3Archiver{uploader: u}
}

// Archive is an implementation of the Archiver interface.
func (a *S3Archiver) Archive(ctx context.Context, logs map[string][]byte) ([]string, error) {
	written := []string{}
	for filename, data := range logs {
		newURL, err := a.uploader.UploadWithContext(ctx, filename, bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		written = append(written, newURL)
	}
	return written, nil
}
