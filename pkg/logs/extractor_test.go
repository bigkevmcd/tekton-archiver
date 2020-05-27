package logs

import (
	"context"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	testPodName = "test-pod-12345"
)

var _ Extractor = (*K8sExtractor)(nil)

func TestPipelineRunLogs(t *testing.T) {
	cl := fake.NewSimpleClientset(makePod())
	e := New(cl)
	e.streamer = func(ns, name string, c kubernetes.Interface) (io.ReadCloser, error) {
		return ioutil.NopCloser(strings.NewReader("testing")), nil
	}
	logs, err := e.PipelineRun(context.Background(), makePipelineRun())
	if err != nil {
		t.Fatal(err)
	}

	want := map[string][]byte{
		testPodName: []byte(`testing`),
	}

	if diff := cmp.Diff(want, logs); diff != "" {
		t.Fatalf("logs don't match:\n%s", diff)
	}
}

func TestPipelineRunLogsPodNotFound(t *testing.T) {
	notFoundErr := apierrors.NewNotFound(schema.GroupResource{Group: "", Resource: "pod"}, testPodName)
	cl := fake.NewSimpleClientset()
	e := New(cl)
	e.streamer = func(ns, name string, c kubernetes.Interface) (io.ReadCloser, error) {
		return nil, notFoundErr
	}

	_, err := e.PipelineRun(context.Background(), makePipelineRun())
	assertErrorMatch(t, "error in opening stream: pod \"test-pod-12345\" not found", err)
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
								PodName: testPodName,
							},
						},
					},
				},
			},
		},
	}
}

func makePod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      testPodName,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "sa",
			RestartPolicy:      corev1.RestartPolicyNever,
			Containers: []corev1.Container{{
				Name:  "nop",
				Image: "nop:latest",
			}},
		},
	}
}

func assertErrorMatch(t *testing.T, s string, e error) {
	t.Helper()
	if s == "" && e == nil {
		return
	}
	if s != "" && e == nil {
		return
	}
	match, err := regexp.MatchString(s, e.Error())
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("failed to match error: %#v with %#v", s, e.Error())
	}
}
