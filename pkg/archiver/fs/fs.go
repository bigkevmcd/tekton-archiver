package archiver

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/spf13/afero"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/bigkevmcd/tekton-archiver/pkg/logs"
)

// FsArchiver implements the LogArchiver interface, and archives to a specific
// filesystem path.
//
// This is useful for testing.
type FsArchiver struct {
	fs       afero.Fs
	basePath string
	extract  logs.Extractor
}

// New creates and returns a new FsArchiver.
func New(fs afero.Fs, e logs.Extractor, basePath string) *FsArchiver {
	return &FsArchiver{fs: fs, extract: e, basePath: basePath}
}

func (a *FsArchiver) ArchivePipelineRun(ctx context.Context, pr *pipelinev1.PipelineRun) ([]string, error) {
	l, err := a.extract.PipelineRun(ctx, pr)
	if err != nil {
		return nil, err
	}
	written := []string{}
	for k, v := range l {
		fname := pathForPR(a.basePath, pr, k+".txt")
		err := afero.WriteFile(a.fs, fname, v, 0644)
		if err != nil {
			return nil, err
		}
		written = append(written, fname)
	}

	return written, nil
}

func (a *FsArchiver) ArchiveTaskRun(ctx context.Context, tr *pipelinev1.TaskRun) ([]string, error) {
	return []string{}, nil
}

func pathForPR(base string, p *pipelinev1.PipelineRun, elements ...string) string {
	return filepath.Join(append([]string{base, p.ObjectMeta.Name}, elements...)...)
}

func metadata(pr *pipelinev1.PipelineRun) ([]byte, error) {
	return json.Marshal(pr)
}
