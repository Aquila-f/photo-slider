package reader

import "context"

type PhotoReader interface {
	ReadPhoto(ctx context.Context, path string) ([]byte, error)
}