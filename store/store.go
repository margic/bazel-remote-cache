package store

import (
	"context"
	"io"
)

// Storer a bazel cache store interface
type Storer interface {
	Get(ctx context.Context, path string, sha string, w io.WriterAt) (int64, error)
	Put(ctx context.Context, path string, sha string, r io.Reader) error
}
