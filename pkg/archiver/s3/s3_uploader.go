package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Uploader wraps the AWS S3 Uploader.
type S3Uploader struct {
	bucketName string
	s3up       *s3manager.Uploader
}

// NewS3Uploader creates and returns a new wrapper around the S3Uploader.
func NewS3Uploader(b string, s3up *s3manager.Uploader) *S3Uploader {
	return &S3Uploader{
		bucketName: b,
		s3up:       s3up,
	}
}

func (u *S3Uploader) UploadWithContext(ctx context.Context, filename string, body io.Reader) (string, error) {
	upParams := &s3manager.UploadInput{
		Bucket: &u.bucketName,
		Key:    &filename,
		Body:   body,
	}
	result, err := u.s3up.UploadWithContext(ctx, upParams)
	if err != nil {
		return "", err
	}
	return result.Location, nil
}
