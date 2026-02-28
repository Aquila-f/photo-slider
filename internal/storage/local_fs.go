package storage

import (
	"context"
	"os"

	"github.com/Aquila-f/photo-slider/internal/domain"
)


type LocalFSProvider struct {
	baseDir string
}


func NewLocalFSProvider(baseDir string) *LocalFSProvider {
	return &LocalFSProvider{baseDir: baseDir}
}

// TODO: consider replacing []domain.FileInfo with iter.Seq2[domain.FileInfo, error]
// for lazy streaming and avoiding full slice allocation on large directories.
func (p *LocalFSProvider) ListDir(ctx context.Context, path string) ([]domain.FileInfo, error) {
	entries, err := os.ReadDir(p.baseDir + "/" + path)
	if err != nil {
		return nil, err
	}
	var files []domain.FileInfo
	for _, e := range entries {
		files = append(files, domain.FileInfo{Name: e.Name(), Path: path, IsDir: e.IsDir()})
	}
	return files, nil
}

func (p *LocalFSProvider) ReadFile(ctx context.Context, filePath string) ([]byte, error) {
	return os.ReadFile(p.baseDir + "/" + filePath)
}