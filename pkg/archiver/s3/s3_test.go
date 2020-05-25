package s3

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/bigkevmcd/tekton-archiver/pkg/logs"
	"github.com/bigkevmcd/tekton-archiver/test"
)

func TestArchivingPipelineRun(t *testing.T) {
	fs := afero.NewMemMapFs()
	a := New(fs, &mockExtractor{}, "/tmp/logs")
	pr := test.MakePipelineRun("my-test-pr", "my-test-pr-pod-12345")

	paths, err := a.ArchivePipelineRun(context.TODO(), pr)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"/tmp/logs/my-test-pr/my-test-pr-pod-12345.txt"}
	if diff := cmp.Diff(want, paths); diff != "" {
		t.Fatalf("failed to archive paths:\n%s", diff)
	}
	b := mustReadFile(t, fs, "/tmp/logs/my-test-pr/my-test-pr-pod-12345.txt")
	if string(b) != "test-output" {
		t.Fatalf("got %s, want %s", string(b), "test-output")
	}
}

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

func mustReadFile(t *testing.T, fs afero.Fs, fname string) []byte {
	t.Helper()
	b, err := afero.ReadFile(fs, fname)
	if err != nil {
		t.Fatalf("failed to read file %s: %s", fname, err)
	}
	return b
}