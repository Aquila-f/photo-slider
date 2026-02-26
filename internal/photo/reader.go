package photo

import (
	"context"
	"os"
)

type Reader interface {
	Read(ctx context.Context, path string) ([]byte, error)
}

type FileReader struct{}

func NewFileReader() *FileReader {
	return &FileReader{}
}

func (r *FileReader) Read(_ context.Context, path string) ([]byte, error) {
	return os.ReadFile(path)
}
