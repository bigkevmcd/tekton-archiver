package s3

import (
	"context"
	"encoding/json"
	"path/filepath"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/bigkevmcd/tekton-archiver/pkg/logs"
)

// "github.com/aws/aws-sdk-go/service/s3/s3manager"
// https://godoc.org/github.com/aws/aws-sdk-go/service/s3/s3manager#NewUploader

// S3Archiver implements the LogArchiver interface, and archives to S3.
type S3Archiver struct {
	uploader Uploader
	bucket   string
	basePath string
	extract  logs.Extractor
}

// New creates and returns a new S3Archiver.
func New(u Uploader, e logs.Extractor, basePath string) *S3Archiver {
	return &S3Archiver{uploader: u, extract: e, basePath: basePath}
}

func (a *S3Archiver) ArchivePipelineRun(ctx context.Context, pr *pipelinev1.PipelineRun) ([]string, error) {
	l, err := a.extract.PipelineRun(ctx, pr)
	if err != nil {
		return nil, err
	}
	written := []string{}
	for k, v := range l {
		fname := pathForPR(a.basePath, pr, k+".txt")
		written = append(written, fname)
	}

	return written, nil
}

func (a *S3Archiver) ArchiveTaskRun(ctx context.Context, tr *pipelinev1.TaskRun) ([]string, error) {
	return []string{}, nil
}

func keyForPR(base string, p *pipelinev1.PipelineRun, elements ...string) string {
	return filepath.Join(append([]string{base, p.ObjectMeta.Name}, elements...)...)
}

func metadata(pr *pipelinev1.PipelineRun) ([]byte, error) {
	return json.Marshal(pr)
}
