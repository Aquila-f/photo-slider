package lister

import "context"

type FileLister interface {
	ListAll(ctx context.Context) ([]string, error)
	GetCompletePath(ctx context.Context, key string) (string, error)
}
