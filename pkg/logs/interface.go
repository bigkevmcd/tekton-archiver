package logs

import (
	"context"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// Extractor defines an interface for getting the logs for a Tekton resource.
type Extractor interface {
	PipelineRun(ctx context.Context, pr *pipelinev1.PipelineRun) (map[string][]byte, error)
	TaskRun(ctx context.Context, pr *pipelinev1.TaskRun) ([]byte, error)
}
