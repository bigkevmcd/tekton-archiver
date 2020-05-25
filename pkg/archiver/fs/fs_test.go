package fs

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"

	"github.com/bigkevmcd/tekton-archiver/pkg/archiver"
)

var _ archiver.Interface = (*FsArchiver)(nil)

func TestArchive(t *testing.T) {
	fs := afero.NewMemMapFs()
	a := New(fs, "/tmp/logs")

	paths, err := a.Archive(context.TODO(), map[string][]byte{
		"my-test-pr/my-test-pr-pod-12345.txt": []byte("test-output"),
		"my-test-pr/my-test-pr-pod-23456.txt": []byte("other-output"),
	})
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"/tmp/logs/my-test-pr/my-test-pr-pod-12345.txt",
		"/tmp/logs/my-test-pr/my-test-pr-pod-23456.txt",
	}
	if diff := cmp.Diff(want, paths); diff != "" {
		t.Fatalf("failed to archive paths:\n%s", diff)
	}
	b := mustReadFile(t, fs, "/tmp/logs/my-test-pr/my-test-pr-pod-12345.txt")
	if string(b) != "test-output" {
		t.Fatalf("got %s, want %s", string(b), "test-output")
	}
	b = mustReadFile(t, fs, "/tmp/logs/my-test-pr/my-test-pr-pod-23456.txt")
	if string(b) != "other-output" {
		t.Fatalf("got %s, want %s", string(b), "other-output")
	}
}

func mustReadFile(t *testing.T, fs afero.Fs, fname string) []byte {
	t.Helper()
	b, err := afero.ReadFile(fs, fname)
	if err != nil {
		t.Fatalf("failed to read file %s: %s", fname, err)
	}
	return b
}
