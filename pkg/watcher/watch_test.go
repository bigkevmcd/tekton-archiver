package watcher

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/bigkevmcd/tekton-archiver/pkg/archiver/fs"
	"github.com/bigkevmcd/tekton-archiver/test"
)

func TestHandlePipelineRun(t *testing.T) {
	mockArchiver := fs.New(afero.NewMemMapFs(), "/tmp")
	logger := zaptest.NewLogger(t, zaptest.Level(zap.WarnLevel))

	err := handlePipelineRun(&mockExtractor{}, mockArchiver, test.MakePipelineRun("test-pr", "test-pr-pod-abc"), logger.Sugar())
	if err != nil {
		t.Fatal(err)
	}
}

type mockExtractor struct{}

func (m *mockExtractor) PipelineRun(ctx context.Context, pr *pipelinev1.PipelineRun) (map[string][]byte, error) {
	return nil, nil
}
