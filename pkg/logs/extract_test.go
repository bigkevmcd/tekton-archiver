package logs

import (
	"context"
	"testing"

	"k8s.io/client-go/kubernetes/fake"
)

func TestPipelineRunLogs(t *testing.T) {
	ctx := context.Background()
	cl := fake.NewSimpleClientset()

	logs, err := PipelineRunLogs(ctx, nil, cl)
	if err != nil {
		t.Fatal(err)
	}

	if logs != nil {
		t.Fatalf("got %#v, want nil", logs)
	}

}
