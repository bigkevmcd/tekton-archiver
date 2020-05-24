package watcher

import (
	"testing"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
)

func TestGetPipelineRunStatus(t *testing.T) {
	statusTests := []struct {
		conditionType   apis.ConditionType
		conditionStatus corev1.ConditionStatus
		want            State
	}{
		{apis.ConditionSucceeded, corev1.ConditionTrue, Successful},
		{apis.ConditionSucceeded, corev1.ConditionUnknown, Pending},
		{apis.ConditionSucceeded, corev1.ConditionFalse, Failed},
	}

	for _, tt := range statusTests {
		w := makePipelineRunWithCondition(tt.conditionType, tt.conditionStatus)
		s := runState(w)
		if s != tt.want {
			t.Errorf("RunState(%s) got %v, want %v", tt.conditionStatus, s, tt.want)
		}
	}
}

func TestTerminalStates(t *testing.T) {
	stateTests := []struct {
		state    State
		terminal bool
	}{
		{Pending, false},
		{Failed, true},
		{Successful, true},
		{Error, true},
		{Unknown, false},
	}

	for _, tt := range stateTests {
		if v := tt.state.Terminal(); v != tt.terminal {
			t.Errorf("Terminal state for %s want %v, got %v", tt.state, tt.terminal, v)
		}
	}

}

func makePipelineRunWithCondition(s apis.ConditionType, c corev1.ConditionStatus) *pipelinev1.PipelineRun {
	return &pipelinev1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pr",
			Namespace: "default",
		},
		Status: pipelinev1.PipelineRunStatus{
			Status: duckv1beta1.Status{
				Conditions: duckv1beta1.Conditions{
					apis.Condition{Type: s, Status: c},
				},
			},
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
