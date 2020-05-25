package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	s3 "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Uploader interface {
	UploadWithContext(ctx aws.Context, input *s3.UploadInput) (*s3.UploadOutput, error)
}
