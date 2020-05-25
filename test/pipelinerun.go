package test

import (
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakePipelineRun(name, podName string) *pipelinev1.PipelineRun {
	return &pipelinev1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Status: pipelinev1.PipelineRunStatus{
			PipelineRunStatusFields: pipelinev1.PipelineRunStatusFields{
				TaskRuns: map[string]*pipelinev1.PipelineRunTaskRunStatus{
					"task1": &pipelinev1.PipelineRunTaskRunStatus{
						Status: &pipelinev1.TaskRunStatus{
							TaskRunStatusFields: pipelinev1.TaskRunStatusFields{
								PodName: podName,
							},
						},
					},
				},
			},
		},
	}
}
