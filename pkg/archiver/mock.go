package archiver

import (
	"context"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/bigkevmcd/tekton-archiver/pkg/logs"
)

var _ logs.Extractor = (*mockExtractor)(nil)

type mockExtractor struct {
}

func (m *mockExtractor) PipelineRun(ctx context.Context, pr *pipelinev1.PipelineRun) (map[string][]byte, error) {
	data := map[string][]byte{
		"my-test-pr-pod-12345": []byte("test-output"),
	}
	return data, nil
}

func (m *mockExtractor) TaskRun(ctx context.Context, pr *pipelinev1.TaskRun) ([]byte, error) {
	return nil, nil
}
