package s3

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bigkevmcd/tekton-archiver/pkg/archiver"
)

var _ archiver.Interface = (*S3Archiver)(nil)

func TestArchive(t *testing.T) {
	uploader := &stubUploader{files: map[string]string{}, bucket: "new-bucket"}
	a := New(uploader)

	paths, err := a.Archive(context.TODO(), map[string][]byte{
		"my-test-pr/my-test-pr-pod-12345.txt": []byte("test-output"),
		"my-test-pr/my-test-pr-pod-23456.txt": []byte("other-output"),
	})
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"new-bucket:my-test-pr/my-test-pr-pod-12345.txt",
		"new-bucket:my-test-pr/my-test-pr-pod-23456.txt",
	}
	if diff := cmp.Diff(want, paths); diff != "" {
		t.Fatalf("failed to archive paths:\n%s", diff)
	}

	wantUploads := map[string]string{
		"new-bucket:my-test-pr/my-test-pr-pod-12345.txt": "test-output",
		"new-bucket:my-test-pr/my-test-pr-pod-23456.txt": "other-output",
	}
	if diff := cmp.Diff(wantUploads, uploader.files); diff != "" {
		t.Fatalf("uploaded files:\n%s", diff)
	}
}

type stubUploader struct {
	bucket string
	files  map[string]string
}

func (s *stubUploader) UploadWithContext(ctx context.Context, filename string, body io.Reader) (string, error) {
	key := fmt.Sprintf("%s:%s", s.bucket, filename)
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	s.files[key] = string(data)
	return key, nil
}
