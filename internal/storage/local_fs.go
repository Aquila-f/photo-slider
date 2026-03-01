package storage

import (
	"context"
	"os"
	"path/filepath"

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
	entries, err := os.ReadDir(filepath.Join(p.baseDir, path))
	if err != nil {
		return nil, err
	}
	var files []domain.FileInfo
	for _, e := range entries {
		files = append(files, domain.FileInfo{Name: e.Name(), Path: filepath.Join(path, e.Name()), IsDir: e.IsDir()})
	}
	return files, nil
}

func (p *LocalFSProvider) Walk(ctx context.Context, root string, maxDepth int) ([]domain.DirSnapshot, error) {
	var snaps []domain.DirSnapshot
	var walk func(dir string, depth int) error
	walk = func(dir string, depth int) error {
		files, err := p.ListDir(ctx, dir)
		if err != nil {
			return err
		}
		snaps = append(snaps, domain.DirSnapshot{Path: dir, Files: files})
		if depth >= maxDepth {
			return nil
		}
		for _, f := range files {
			if f.IsDir {
				if err := walk(filepath.Join(dir, f.Name), depth+1); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return snaps, walk(root, 0)
}

func (p *LocalFSProvider) ReadFile(ctx context.Context, filePath string) ([]byte, error) {
	return os.ReadFile(filepath.Join(p.baseDir, filePath))
}
