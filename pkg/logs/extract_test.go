package logs

import (
	"context"
	"testing"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestPipelineRunLogs(t *testing.T) {
	ctx := context.Background()
	// https://godoc.org/k8s.io/api/core/v1#Pod
	cl := fake.NewSimpleClientset()

	logs, err := PipelineRunLogs(ctx, makePipelineRun(), cl)
	if err != nil {
		t.Fatal(err)
	}

	if logs != nil {
		t.Fatalf("got %#v, want nil", logs)
	}
}

func makePipelineRun() *pipelinev1.PipelineRun {
	return &pipelinev1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pr",
			Namespace: "default",
		},
		Status: pipelinev1.PipelineRunStatus{
			PipelineRunStatusFields: pipelinev1.PipelineRunStatusFields{
				TaskRuns: map[string]*pipelinev1.PipelineRunTaskRunStatus{
					"task1": &pipelinev1.PipelineRunTaskRunStatus{
						Status: &pipelinev1.TaskRunStatus{
							TaskRunStatusFields: pipelinev1.TaskRunStatusFields{
								PodName: "test-pod",
							},
						},
					},
				},
			},
		},
	}
}
