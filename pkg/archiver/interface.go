package archiver

import (
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// LogArchiver provides the core interface for archiving the output from
// PipelineRuns.
type LogArchiver interface {
	// Archive the PipelineRun output and return a URL to retrieve
	// the contents later, or an error.
	ArchivePipelineRun(*pipelinev1.PipelineRun, []byte) (string, error)
}
