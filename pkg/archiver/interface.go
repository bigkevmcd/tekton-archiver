package archiver

import "context"

// Interface implementations persist a map of files to a store.
type Interface interface {
	Archive(context.Context, map[string][]byte) ([]string, error)
}
