package fs

import (
	"context"
	"path/filepath"

	"github.com/spf13/afero"
)

// FsArchiver implements the Archiver interface, and archives to a specific
// filesystem path.
//
// This is useful for testing.
type FsArchiver struct {
	fs       afero.Fs
	basePath string
}

// New creates and returns a new FsArchiver.
func New(fs afero.Fs, basePath string) *FsArchiver {
	return &FsArchiver{fs: fs, basePath: basePath}
}

func (a *FsArchiver) Archive(ctx context.Context, logs map[string][]byte) ([]string, error) {
	written := []string{}
	for k, v := range logs {
		fname := filepath.Join(a.basePath, k)
		err := afero.WriteFile(a.fs, fname, v, 0644)
		if err != nil {
			return nil, err
		}
		written = append(written, fname)
	}
	return written, nil
}
